.PHONY: docker
docker:
	@docker build -t ahyang98/vote_worker ./vote/worker
	@docker build -t ahyang98/vote_vote ./vote/vote
	@docker build -t ahyang98/vote_seed_data ./vote/seed-data
	@docker build -t ahyang98/vote_result ./vote/result
.PHONY: helm
helm:
	@helm upgrade --install vote ./chart/ -n vote --create-namespace