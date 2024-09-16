package database

var (
	CollectTenderFromRow = "id, version, name, description, status, type, organization_id, creator_username, created_at"
	InsertTender         = "INSERT INTO Tenders (version, name, description, status, type, organization_id, creator_username)"
	SelectTender         = "SELECT id, version, name, description, status, type, organization_id, creator_username, created_at FROM Tenders"
	InsertTenderVersion  = "INSERT INTO Tenders (id, version, name, description, status, type, organization_id, creator_username)"
	InsertBid            = "INSERT INTO Bids (tender_id, name, description, creator_username, organization_id)"
	SelectBid            = "SELECT id, tender_id, name, description, status, creator_username, organization_id FROM Bids"
	CollectBidFromRow    = "id, tender_id, name, description, status, creator_username, organization_id, created_at"
)
