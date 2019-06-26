module github.com/phillbaker/terraform-provider-elasticsearch

go 1.12

require (
	git.apache.org/thrift.git v0.12.0 // indirect
	github.com/aws/aws-sdk-go v1.19.38
	github.com/deoxxa/aws_signing_client v0.0.0-20161109131055-c20ee106809e
	github.com/grpc-ecosystem/grpc-gateway v1.6.2 // indirect
	github.com/hashicorp/terraform v0.12.0
	github.com/olivere/elastic v6.2.18+incompatible
	github.com/olivere/elastic/v7 v7.0.2-0.20190606091611-4dacbebcb82a
	gopkg.in/olivere/elastic.v5 v5.0.81
	gopkg.in/olivere/elastic.v6 v6.2.19-0.20190606093138-f7db55b7060d
)

replace github.com/olivere/elastic/v7 => github.com/jfroche/elastic/v7 v7.0.4-0.20190626230911-ba08bbca38c6
