package elastic

const MAPPING = `
{
	"settings": {
		"analysis": {
			"tokenizer": {
			  "trigram_tokenizer": {
				"type": "ngram",
				"min_gram": 3,
				"max_gram": "3",
				"token_chars": [
				  "letter",
				  "digit"
				]
			  }
			},
			"analyzer": {
				"trigrams": {
					"type": "custom",
					"tokenizer": "trigram_tokenizer",
					"filter": [
						"lowercase"
					]
				}
			}
		}
	}, 
	"mappings": {
		"properties": {
			"username": {
				"type": "text",
				"analyzer": "trigrams"
			},
			"ID": {
				"type": "keyword"
			},
			"has_picture": {
				"type": "boolean"
			}
		}
	}
}`
