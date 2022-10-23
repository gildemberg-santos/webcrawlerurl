# WebCrawlerUrl
Esse serviço e composto por duas funções do Google Cloud Function, para mapear as urls de uma página e consultar esses dados em um banco de dados MongoDB
## Function #1
### WebCrawlerUrlHttp
Responsavel por pegar os dados no banco de dados

#### Instalação Google Cloud Function
* Crie uma arquivos go.mod na raiz com todos os dados do go.mod local

* Copie e cole todos os arquivos da pasta p para a raiz do GCF.

* Configure o nome de metodo incial para WebCrawlerUrlHttp

* Crie uma variavel de ambiente chamada PROJECT_ID com o nome do seu projeto Google Cloud.

* Crie uma variavel de ambiente chamada TOPIC_NAME com o nome do topico do seu pubsub.

## Function #2
### WebCrawlerUrlPubSub
Vai ser responsavel por pesquisar todos os links e gravar em uma banco de dados

#### Instalação Google Cloud Function
* Crie uma arquivos go.mod na raiz com todos os dados do go.mod local

* Copie e cole todos os arquivos da pasta p para a raiz do GCF.

* Configure o nome de metodo incial para WebCrawlerUrlPubSub

* Crie uma variavel de ambiente chamada MONGO_STR_CONNECTION com sua string de conecxão com seu  MongoDB.

