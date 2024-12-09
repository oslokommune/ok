tape: pre-steps tape.vhs
	vhs tape.vhs
	open tape.webm
	$(MAKE) post-steps

version: pre-steps version.vhs
	vhs version.vhs
	open version.webm
	$(MAKE) post-steps

interactive: pre-steps interactive.vhs
	vhs interactive.vhs
	open interactive.webm
	$(MAKE) post-steps

pre-steps:
	cp ../../test/terraform/packages.yml .
	cp -r ../../test/terraform/config .

post-steps:
	rm packages.yml
	rm -rf config
