# Changelog

## v1.6.0
* feat: support the user interface to grant or revoke privileges
* feat: support to sort the data by the sortRules firstly, when quering documents
* feat: return [FailIndexedDocumentSets] in describing collectionView's response
* feat: return the number of collections/collectionViews in listing database's response 

## v1.5.5
* refactor: update default stopwords file

## v1.5.4
* feat: support search radius for searching in collection
 
## v1.5.3
* fix: optimize warning message for modifing vector indexes interface

## v1.5.2
* fix: optimize the interface modifing vector indexes 

## v1.5.1
* fix: replace large file-sized dictionaries in Git LFS with URLs to enable the usage of NewBM25Encoder with language-specific dictionaries when customers refer to the "tcvdbtext" directory in this Go package.

## v1.5.0
* feat: support binary vector
* feat: support to return documents' count when using Count function
* feat: support to delete documents with limit
* feat: add params terminateAfter and cutoffFrequency, when hybrid search
* feat: support to hybrid search by text in embedding collection
* feat: support to modify vector indexes

## v1.4.9
* 为tcvectordb包的部分接口，增加注释
* ExistsDatabase实现方法中：增加对AI_DB类型数据库的判断

## v1.4.8
* CreateCollectionView和LoadAndSplitText接口支持设置文件解析的parsingType，可设置为AlgorithmParsing/VisionModelParsing

## v1.4.7
* 移除buildtag，支持windows平台上编译运行

## v1.4.6
* 支持BAAI/bge-m3新模型，也支持创建embedding collection时使用string直接设置模型

## v1.4.5
* 更换依赖cgo的分词包，为纯go实现的分词包，以更好的支持跨平台编译

## v1.4.4
* 新增/index/add接口实现

## v1.0.0

### DatabaseInterface
* DropDatabase添加返回值DropDatabaseResult
* ListDatabase返回值修改为ListDatabaseResult结构，将原[]Database放入result中
* 新增CreateAIDatabase接口，创建AI的database
* 新增DropAIDatabase接口，删除AI的database
* 新增AIDatabase接口，用于非API调用获取aidb对象

### CollectionInterface
* DropCollection添加返回值DropCollectionResult
* TruncateCollection返回值AffectedCount修改为TruncateCollectionResult
* ListCollection返回值[]*Collection修改为ListCollectionResult

### AICollectionViewInterface
* 新增AI collectionView相关接口

### AliasInterface
* 接口名称修改，AliasSet替换为SetAlias，返回值修改为SetAliasResult
* 接口名称修改，AliasDelete替换为DeleteAlias， 返回值修改为DeleteAliasResult

### AIAliasInterface
* 新增AI collection的别名相关接口

### IndexInterface
* IndexRebuild修改名称为RebuildIndex, 返回值修改为RebuildIndexResult

### DocumentInterface
* Upsert移动buildIndex参数到option中；返回值添加UpsertDocumentResult
* Query移动retrieveVector到option中，返回值修改为QueryDocumentResult
* Search\SearchById\SearchByText移动filter、hnswparam、retrieveVector、limit到option中；返回值修改为SearchDocumentResult
* Delete移动documentIds到option中，返回值添加DeleteDocumentResult

### AIDocumentSetsInterface
* 新增AI documentSet相关接口