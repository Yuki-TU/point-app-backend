# makeを打った時のコマンド
.DEFAULT_GOAL := help

# 環境変数
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
IMAGE_TAG=$(git rev-parse HEAD)
ECR_REGISTRY=${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com

# データベース
DB_HOST=db
DB_PORT=3306
DB_USER=admin
DB_PASSWORD=password
DB_NAME=point_app
DB_TEST_NAME=point_app_test

.PHONY: build
build: ## デプロイ用のDockerイメージをビルド
	docker image build \
		-t ${ECR_REGISTRY}/point-app-backend:latest \
		-t ${ECR_REGISTRY}/point-app-backend:${IMAGE_TAG}  \
		--target deploy ./

.PHONY: build-up
build-up: ## Dockerイメージをビルドしてコンテナを起動
	docker compose up -d --build

.PHONY: push
push: ## ECRにプッシュ
	aws ecr --region ap-northeast-1 get-login-password | docker login --username AWS --password-stdin https://${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/point-app-backend
	docker image push -a ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/point-app-backend

.PHONY: in
in: ## アプリケーションのコンテナに入る（ホスト）
	docker compose exec app sh

.PHONY: up
up: ## ホットリロード付きでdocker compose upを実行（ホスト）
	docker compose up -d app db cache aws

.PHONY: down
down: ## docker compose downを実行（ホスト）
	@docker compose down

.PHONY: format
log: ## Docker composeのログを表示（ホスト）
	@docker compose logs app -f

.PHONY: ps
ps: ## コンテナの状態を確認（ホスト）
	docker compose ps

.PHONY: rsa 
rsa: down build-up ## 全てのコンテナを削除して、ビルドして、起動

.PHONY: dry-migrate
dry-migrate: ## マイグレーションの試行（マイグレーション時に発行されるDDL確認）
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} -P ${DB_PORT} ${DB_NAME} --dry-run < ./_tools/mysql/schema.sql

.PHONY: migrate
migrate:  ## マイグレーションを実行（コンテナ）
	@mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} -P ${DB_PORT} ${DB_NAME} < ./_tools/mysql/schema.sql
	@if [ ${GO_ENV} == development ]; then \
		mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} -P ${DB_PORT} ${DB_TEST_NAME} < ./_tools/mysql/schema.sql; \
	fi

.PHONY: seed
seed: ## データ挿入（コンテナ）
	mariadb --skip-ssl ${DB_NAME} -h ${DB_HOST} -u ${DB_USER} -p${DB_PASSWORD} < ./_tools/mysql/seed.sql 

model: ## モデル作成
	rm -rf ./repository/entities
	mkdir -p ./repository/entities
	xo schema 'mysql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}/${DB_NAME}?parseTime=true&sql_mode=ansi' -o ./repository/entities --go-field-tag='json:"{{ .SQLName }}" db:"{{ .SQLName }}"'

.PHONY: rdm
rdm: ## 送信メールを確認
	@if [ ${CONTAINER_ENV} ]; then \
		curl -v http://aws:4566/_aws/ses/ | jq . | tail -n 18 | head -n 16; \
	else \
		curl -v http://localhost:4566/_aws/ses/ | jq . | tail -n 18 | head -n 16; \
	fi

.PHONY: create-key
create-key: ## JWTに必要なキーを作成
	openssl genrsa 4096 > ./auth/certificate/secret.pem
	openssl rsa -pubout < ./auth/certificate/secret.pem > ./auth/certificate/public.pem

.PHONY: format
format: ## コードフォーマット
	@echo "フォーマット"
	@if [ ${CONTAINER_ENV} ]; then \
		gofmt -l -s -w .; \
		goimports -w -l .; \
	else \
		docker compose exec app gofmt -l -s -w .; \
		docker compose exec app goimports -w -l .; \
	fi


.PHONY: lint
lint: format ## リンター(golangci-lint)
	@echo "リンター"
	@if [ ${CONTAINER_ENV} ]; then \
		golangci-lint run; \
	else \
		docker compose exec app golangci-lint run; \
	fi


.PHONY: moq
moq: ## モック作成(コンテナ内)
	# サービスのモック生成中
	@docker compose exec app moq -fmt goimports -out ./handler/moq_test.go ./handler \
					RegisterUserService \
					RegisterTemporaryUserService \
					SigninService \
					GetUsersService \
					UpdatePasswordService \
					UpdateAccountService \
					ResetPasswordService \
					SendPointService \
					SignoutService \
					GetAccountService \
					UpdateTemporaryEmailService \
					GetNotificationService \
					GetNotificationsService \
					GetUncheckedNotificationCountService \
					DeleteUserService

	# リポジトリのモック生成中
	@docker compose exec app moq -fmt goimports -out ./service/moq_test.go -skip-ensure -pkg service ./domain \
					Cache \
					TokenGenerator \
					UserRepo \
					PointRepo \
					NotificationRepo
	@docker compose exec app moq -fmt goimports -out ./service/repogitory_moq_test.go -skip-ensure -pkg service ./repository Beginner Preparer Execer Queryer

.PHONY: mock
mock: ## モック作成
	mockgen -source=./batch/controller/usecase.go -destination=./batch/controller/_mock/mock_usecase.go
	mockgen -source=./repository/repository.go -destination=./repository/_mock/mock_repository.go
	mockgen -source=./domain/interface.go -destination=./domain/_mock/mock_interface.go

.PHONY: test
test: ## テスト実行
	# テスト実行中
	@if [ ${CONTAINER_ENV} ]; then \
		go test -cover -race -shuffle=on ./...; \
	else \
		docker compose exec app go test -cover -race -shuffle=on ./...; \
	fi

.PHONY: coverage
coverage: ## カバレッジファイル作成・表示（ホスト側）
	# テスト実行中
	@docker compose exec app go test -cover ./... -coverprofile=cover.out
	# HTMLに変換中
	@docker compose exec app go tool cover -html=cover.out -o tmp/cover.html
	@docker compose exec app rm cover.out
	# ブラウザ表示
	@open ./tmp/cover.html

.PHONY: wire
wire: ## API用のDIファイル生成
	@wire ./router

.PHONY: wire-b
wire-b: ## バッチ用のDIファイル生成
	@wire ./batch/wire

.PHONY: batch
batch: ## バッチ用アプリケーションのビルド
	@if [ ${CONTAINER_ENV} ]; then \
		cd /app/cmd/batch && go install; \
	else \
		docker compose exec app sh -c "cd /app/cmd/batch && go install"; \
	fi

.PHONY: db
db: ## データベースに入る
	@docker compose exec db mysql ${DB_NAME}

.PHONY: env
env: ## 環境変数を表示
	@cp .secrets.example .secrets

.PHONY: act
act: ## テストを実行
	act pull_request --secret-file .secrets --container-architecture linux/amd64

.PHONY: help
help: ## オプションを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
