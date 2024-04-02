module metadata_management

replace file_manager => ../file_manager

replace log_manager => ../log_manager

replace buffer_manager => ../buffer_manager

replace tx => ../tx

replace record_manager => ../record_manager

replace query => ../query

replace index_manager => ../index_manager

go 1.19

require (
	buffer_manager v0.0.0-00010101000000-000000000000
	file_manager v0.0.0-00010101000000-000000000000
	index_manager v0.0.0-00010101000000-000000000000
	log_manager v0.0.0-00010101000000-000000000000
	query v0.0.0-00010101000000-000000000000
	record_manager v0.0.0-00010101000000-000000000000
	tx v0.0.0-00010101000000-000000000000
)
