package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	utils "github.com/lakshay88/reward-management-system/Utils"
	"github.com/lakshay88/reward-management-system/config"
	"github.com/lakshay88/reward-management-system/database/models"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	connection *sql.DB
}

func ConnectionToPostgres(cfg config.DatabaseConfig) (Database, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	connection, err := sql.Open(cfg.Driver, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Set connection limits
	connection.SetMaxOpenConns(25)
	connection.SetMaxIdleConns(25)
	connection.SetConnMaxLifetime(5 * time.Minute)

	if err := connection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{connection: connection}, nil
}

func (db *PostgresDB) Close() error {
	return db.connection.Close()
}

func (db *PostgresDB) CreateUser(user *models.User) (*models.User, error) {
	// Insert user into the database
	query := `INSERT INTO users (username, email, user_password) VALUES ($1, $2, $3) RETURNING id, created_on`

	// Execute the query
	err := db.connection.QueryRow(query, user.Username, user.Email, user.UserPassword).Scan(&user.ID, &user.CreatedOn)
	if err != nil {
		return nil, err
	}

	// empty user
	user.UserPassword = ""

	return user, nil
}

func (db *PostgresDB) GetUserByID(userId int, user *models.User) (*models.User, error) {

	// Handling Nil
	if user == nil {
		user = &models.User{}
	}

	query := `SELECT id, username, email, created_on FROM users WHERE id = $1`
	err := db.connection.QueryRow(query, userId).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedOn)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User with ID %s not found", userId)
	} else if err != nil {
		return nil, fmt.Errorf("Failed to fetch user: %v", err)
	}

	return user, nil
}

func (db *PostgresDB) GetUserByEmail(userEmail string, user *models.User) (*models.User, error) {

	// Handling Nil
	if user == nil {
		user = &models.User{}
	}

	query := `SELECT id, username, email, user_password, created_on FROM users WHERE email = $1`
	err := db.connection.QueryRow(query, userEmail).Scan(&user.ID, &user.Username, &user.Email, &user.UserPassword, &user.CreatedOn)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User with Email %s not found", userEmail)
	} else if err != nil {
		return nil, fmt.Errorf("Failed to fetch user: %v", err)
	}

	return user, nil
}

func (db *PostgresDB) PointBalance(userId string) error {

	return nil
}

func (db *PostgresDB) AddTransaction(txn *models.Transaction) (*models.Transaction, error) {

	// Handling Nil
	if txn == nil {
		txn = &models.Transaction{}
	}

	// User validation
	var userCount int
	checkUserQuery := `SELECT COUNT(1) FROM users WHERE id = $1`
	err := db.connection.QueryRow(checkUserQuery, txn.UserID).Scan(&userCount)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if user exists: %v", err)
	}

	if userCount == 0 {
		return nil, fmt.Errorf("User with ID %d does not exist", txn.UserID)
	}

	txn.TransactionID = uuid.New().String()
	if txn.TransactionDate.IsZero() {
		txn.TransactionDate = time.Now()
	}

	// Points Calculation
	var pointsEarned int
	categoryMultiplier := utils.GetCategoryMultiplier(txn.Category)
	pointsEarned = int(txn.TransactionAmount) * categoryMultiplier

	// Transaction add logic
	transactionQuery := `INSERT INTO transactions (transaction_id, user_id, transaction_amount, category, transaction_date, product_code, points_earned) 
						VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var transactionID int
	err = db.connection.QueryRow(transactionQuery, txn.TransactionID, txn.UserID, txn.TransactionAmount,
		txn.Category, txn.TransactionDate, txn.ProductCode, pointsEarned).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert transaction: %v", err)
	}

	// Update if points_balance already exist
	pointsBalanceQuery := `UPDATE points_balance SET total_points = total_points + $1 WHERE user_id = $2`
	result, err := db.connection.Exec(pointsBalanceQuery, pointsEarned, txn.UserID)
	if err != nil {
		return nil, fmt.Errorf("Failed to update points balance: %v", err)
	}

	// If no rows were affected create a new entry.
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		pointsBalanceQuery := `INSERT INTO points_balance (user_id, total_points) VALUES ($1, $2)`
		_, err = db.connection.Exec(pointsBalanceQuery, txn.UserID, pointsEarned)
		if err != nil {
			return nil, fmt.Errorf("Failed to create points balance: %v", err)
		}
	}

	err = db.LogPointsHistory(txn.UserID, pointsEarned, "earn", "Points earned for transaction")
	if err != nil {
		return nil, fmt.Errorf("Failed to log points history: %v", err)
	}

	txn.ID = transactionID
	txn.PointsEarned = pointsEarned
	return txn, nil
}

func (db *PostgresDB) GetPointsBalance(userID int) (models.PointsBalance, error) {
	var balance models.PointsBalance
	query := `SELECT total_points, points_redeemed FROM points_balance WHERE user_id = $1`
	err := db.connection.QueryRow(query, userID).Scan(&balance.TotalPoints, &balance.PointsRedeemed)
	if err == sql.ErrNoRows {
		return balance, fmt.Errorf("User with ID %d has no points balance", userID)
	} else if err != nil {
		return balance, err
	}
	return balance, nil
}

func (db *PostgresDB) GetPointsHistory(userID, page, limit int, startDate, endDate, transactionType string) ([]models.PointsHistory, error) {
	offset := (page - 1) * limit

	query := `SELECT points, points_type, reason, date FROM points_history 
              WHERE user_id = $1`
	args := []interface{}{userID}

	if startDate != "" {
		query += " AND date >= $" + fmt.Sprint(len(args)+1)
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND date <= $" + fmt.Sprint(len(args)+1)
		args = append(args, endDate)
	}
	if transactionType != "" {
		query += " AND points_type = $" + fmt.Sprint(len(args)+1)
		args = append(args, transactionType)
	}

	query += " ORDER BY date DESC LIMIT $" + fmt.Sprint(len(args)+1) + " OFFSET $" + fmt.Sprint(len(args)+2)
	args = append(args, limit, offset)

	rows, err := db.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PointsHistory
	for rows.Next() {
		var entry models.PointsHistory
		if err := rows.Scan(&entry.Points, &entry.PointsType, &entry.Reason, &entry.Date); err != nil {
			return nil, err
		}
		history = append(history, entry)
	}

	return history, nil
}

func (db *PostgresDB) GetAvailablePoints(userID int) (int, error) {
	var totalPoints int
	query := `
		SELECT COALESCE(SUM(points), 0) 
		FROM points_history 
		WHERE user_id = $1 AND points_type = 'earn' AND date >= NOW() - INTERVAL '1 year'`
	err := db.connection.QueryRow(query, userID).Scan(&totalPoints)
	if err != nil {
		return 0, fmt.Errorf("Failed to retrieve available points: %v", err)
	}
	return totalPoints, nil
}

func (db *PostgresDB) DeductPoints(userID int, pointsToRedeem int) (int, error) {
	// Update points balance
	updateBalanceQuery := `
		UPDATE points_balance 
		SET total_points = total_points - $1, points_redeemed = points_redeemed + $1 
		WHERE user_id = $2 RETURNING total_points`
	var remainingBalance int
	err := db.connection.QueryRow(updateBalanceQuery, pointsToRedeem, userID).Scan(&remainingBalance)
	if err != nil {
		return 0, fmt.Errorf("Failed to update points balance: %v", err)
	}
	return remainingBalance, nil
}

func (db *PostgresDB) LogPointsHistory(userID int, points int, pointsType string, reason string) error {
	pointsHistoryQuery := `
		INSERT INTO points_history (user_id, points, points_type, reason, date)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := db.connection.Exec(pointsHistoryQuery, userID, points, pointsType, reason, time.Now())
	if err != nil {
		return fmt.Errorf("Failed to log points history: %v", err)
	}
	return nil
}

func (db *PostgresDB) TransactionOlderThanGivenTime(timePeriod time.Time) ([]models.Transaction, error) {
	// Prepare a slice to store the transactions
	var transactions []models.Transaction

	rows, err := db.connection.Query(`
		SELECT user_id, transaction_id, points_earned, transaction_date 
		FROM transactions
		WHERE transaction_date <= $1 AND points_earned > 0
	`, timePeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var txn models.Transaction
		err := rows.Scan(&txn.UserID, &txn.TransactionID, &txn.PointsEarned, &txn.TransactionDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}

		transactions = append(transactions, txn)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error : %v", err)
	}

	return transactions, nil
}

func (db *PostgresDB) ExpirePoints(userId int, transactionID string, pointsEarned int, transactionDate time.Time) error {
	_, err := db.connection.Exec(`
	UPDATE points_balance 
	SET total_points = total_points - $1 
	WHERE user_id = $2
`, pointsEarned, userId)
	if err != nil {
		return fmt.Errorf("failed to update points balance: %v", err)
	}

	err = db.LogPointsHistory(userId, pointsEarned, "expired", "Point expired due to inactivity")
	if err != nil {
		return fmt.Errorf("Failed to log points history: %v", err)
	}

	fmt.Printf("Expired %d points for user %d, transaction %s.\n", pointsEarned, userId, transactionID)
	return nil
}
