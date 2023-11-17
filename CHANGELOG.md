# Changelog
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

### AICollectionInterface
* 新增AI collection相关接口

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

### AIDocumentInterface
* 新增AI document相关接口