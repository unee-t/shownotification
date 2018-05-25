aws --profile uneet-dev sns publish \
	--topic-arn arn:aws:sns:ap-southeast-1:812644853088:atest \
	--message "Hello World $(date)"
