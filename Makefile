.PHONY: test
test:
	@pipenv run test

.PHONY: vuln
vuln:
	@pipenv check

.PHONY: init
init:
	@pipenv install --dev
