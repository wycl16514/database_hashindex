module query

replace file_manager => ../file_manager

replace log_manager => ../log_manager

replace buffer_manager => ../buffer_manager

replace tx => ../tx

replace record_manager => ../record_manager

replace parser => ../parser

replace query => ../query

go 1.19

require (
	buffer_manager v0.0.0-00010101000000-000000000000
	file_manager v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.1
	log_manager v0.0.0-00010101000000-000000000000
	record_manager v0.0.0-00010101000000-000000000000
	tx v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
