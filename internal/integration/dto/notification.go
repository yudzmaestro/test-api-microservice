package dto

type NotificationRequestDTO struct {
	NotifyID		string		`json:"notify_id"`
	LoanID			string		`json:"loan_id"`
	CustomerID		string		`json:"customer_id"`
	Type			int			`json:"type"`
	Title			string		`json:"title"`
	Content			string		`json:"content"`
	Link			string		`json:"link"`
	Datetime		string		`json:"datetime"`
	Signature		string		`json:"signature"`
}

type NotificationResponseDTO struct {
	Code		int			`json:"code"`
	Message		string		`json:"message"`
}