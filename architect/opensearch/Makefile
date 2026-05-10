

local_start: docker-compose up

mock_gen_repo:
			mockgen -source=repository/es_indexer_repository.go -destination=mocks/mock_es_indexer_repository.go -package=mocks &&
			mockgen -source=repository/es_bulk_indices_repository.go -destination=mocks/mock_es_bulk_indices_repository -package=mocks &&
			mockgen -source=controller/indexer_controller.go -destination=mocks/mock_indexer_controller.go -package=mocks &&
			mockgen -source=controller/search_controller.go -destination=mocks/mock_search_controller.go -package=mocks

