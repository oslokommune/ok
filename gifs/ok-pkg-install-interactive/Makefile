main: pre-steps main.vhs
	vhs main.vhs
	open ok-pkg-install-interactive.webm
	$(MAKE) post-steps

version: pre-steps version.tape
	vhs version.tape
	open version.webm
	$(MAKE) post-steps

interactive: pre-steps interactive.tape
	vhs interactive.tape
	open interactive.webm
	$(MAKE) post-steps

pre-steps:
	cp ../../test/terraform/packages.yml .
	cp -r ../../test/terraform/config .

post-steps:
	rm packages.yml
	rm -rf config
