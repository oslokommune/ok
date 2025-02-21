# Test data structure

In order to understand what the command tests do, it is useful to look at the test data structure.

```shell
install/ # This is a command to test
    recursive/ # This is a test case
        expected/ # This is the expected output
        input/ # Files to be used as input for the test case.
            
            # Files in this dir will used as working directory for ok commands.
            # These files are copied to a temporary directory before running ok commands there. 
            root/
                app-hello/ # Contains a rendered Boilerplate template.

update/ # Another command to test. Folows the same structure as for the "install" folder.

```
