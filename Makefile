.PHONY: test
test:
	@pipenv run test

.PHONY: init
init:
	@pipenv install --dev
