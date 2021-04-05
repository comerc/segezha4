build-image:
	docker pull chromedp/headless-shell:latest
	docker build -fDockerfile -tsegezha4 .

up:
	docker-compose -fdocker-compose.yml up

run-headless:
	docker run -d -p 9222:9222 --rm --name headless-shell chromedp/headless-shell
