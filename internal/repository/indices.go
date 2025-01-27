package repository

// Index mappings
var indexMappings = map[string]string{
	"users": `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"email": { "type": "keyword" },
				"password": { "type": "keyword" },
				"name": { "type": "text" },
				"role": { "type": "keyword" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" },
				"deleted_at": { "type": "date" }
			}
		}
	}`,
	"social_accounts": `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"user_id": { "type": "keyword" },
				"platform": { "type": "keyword" },
				"access_token": { "type": "keyword" },
				"refresh_token": { "type": "keyword" },
				"account_name": { "type": "keyword" },
				"status": { "type": "keyword" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" },
				"deleted_at": { "type": "date" }
			}
		}
	}`,
	"posts": `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"user_id": { "type": "keyword" },
				"content": { "type": "text" },
				"media_urls": { "type": "keyword" },
				"platforms": { "type": "keyword" },
				"status": { "type": "keyword" },
				"scheduled_for": { "type": "date" },
				"published_at": { "type": "date" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" },
				"deleted_at": { "type": "date" }
			}
		}
	}`,
	"media": `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"user_id": { "type": "keyword" },
				"url": { "type": "keyword" },
				"type": { "type": "keyword" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" },
				"deleted_at": { "type": "date" }
			}
		}
	}`,
	"analytics": `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"post_id": { "type": "keyword" },
				"platform": { "type": "keyword" },
				"likes": { "type": "integer" },
				"comments": { "type": "integer" },
				"shares": { "type": "integer" },
				"reach": { "type": "integer" },
				"engagement": { "type": "float" },
				"recorded_at": { "type": "date" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" },
				"deleted_at": { "type": "date" }
			}
		}
	}`,
}
