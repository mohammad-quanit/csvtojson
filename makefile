tests:
	go test -v

test_func:
	@echo "Start Testing '$(FUNC)' function"
	go test -run $(FUNC) -v

git_done:
	git add .
	git commit -m "$(MSG)"
	git push
