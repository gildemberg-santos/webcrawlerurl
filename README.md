# WebCrawlerUrl
Esse serviço e composto por duas funções do Google Cloud Function, para mapear as urls de uma página e consultar esses dados em um banco de dados MongoDB
## Function #1
### WebCrawlerUrlHttp
Responsavel por pegar os dados no banco de dados e enviar a Company e Link do site para uma fila do PubSub.

#### Instalação Google Cloud Function
* Crie uma arquivos go.mod na raiz com todos os dados do go.mod local

* Copie e cole todos os arquivos da pasta p para a raiz do GCF.

* Configure o nome de metodo incial para WebCrawlerUrlHttp

* Crie uma variavel de ambiente chamada **GOOGLE_CLOUD_PROJECT** com o nome do seu projeto Google Cloud.

* Crie uma variavel de ambiente chamada **GOOGLE_TOPIC_NAME** com o nome do topico do seu pubsub.

. Crir uma variavel de ambiente chamada **MONGO_STR_CONNECTION** com sua string de conecxão com seu MongoDB

## Function #2
### WebCrawlerUrlPubSub
Vai ser responsavel por pesquisar todos os links e gravar em uma banco de dados

#### Instalação Google Cloud Function
* Crie uma arquivos **go.mod** na raiz com todos os dados do **go.mod** local

* Copie e cole todos os arquivos da pasta p para a raiz do GCF.

* Configure o nome de metodo incial para WebCrawlerUrlPubSub

* Crie uma variavel de ambiente chamada **MONGO_STR_CONNECTION** com sua string de conecxão com seu MongoDB valor padrão mongodb://localhost:27017/.

* Crie uma variavel de ambiente chamada **LOGS_SHOW** para exibir os logs valores aceitos 0 e 1, valor padrão e 0.

* Crie uma variavel de ambiente chamada **LINKS_MAX** e a quantidade maxima de links simultaneos valor aceitos são inteiros positivos, valor padrão e 5.

* Crie uma variavel de ambiente chamada **LOOP_MAX** e a quantidade de blocos de links simultaneos, valor aceitos são inteiros positicos, valor padrão e 1.

## MongoDB

### Configurar collection
Nome do banco **webcrawler**, nome da collection **links**

Criar Index de campos unicos para a company e o link
``` MongoDB
use.webcrawler
db.links.createIndex({company: 1, link: 1}, {unique: true})
```
