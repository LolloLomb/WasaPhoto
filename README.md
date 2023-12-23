# WasaPhoto
Share your flash!   

## v2.0.4   
Utils added   
All responses has been added   
Fixed rows.Err() checks   
Fixed comments   
TODO: Prepare DB for testing   
TODO: Test responses   

## v2.0.3   
Authorization Implemented   

## v2.0.2 Changelog   
Implemented all methods   
Introduced /profile path   
Fixed unused components warnings   
Fixed rows.Err() checks   
TODO: Testing   
TODO: Implement Authorization   

## v2.0.1 Changelog
Introduced binary data format for images   
Added controls over banUser and followUser methods   
Introduced some utils for db operations   

## v2.0.0 Changelog
Changed security scheme application   
Changed application/json request and response format ( <type> ---> object )   
Implemented first simple methods   

## v1.0.4 Changelog
Fixed schemas properties types  
Fixed identifier behaviour (string --> integer)  

## v1.0.3 Changelog
Fixed date-time minLength maxLength len(example) mismatch (19 --> 20)   

## v1.0.2 Changelog
Removed addFollower and removeFollower operation   
User identifiers parameters now reference the same uid parameter structure   
Removed getFollowingId   
Delete operation now has 204 instead of 200 as success message   
Some operation now use put method instead of post   
Fixed a typo in setMyUserName opertionId   
Responses schemas are now simpler   
Fixed 200 response from getUserProfile   
Fixed date-time pattern, now it follows example like "2017-07-21T17:32:28Z"   
Fixed reference in methods' requestBody   
   
## v1.0.1 Changelog
Added Tags Array  
Common path parameters fixed  
Defined pattern, minLength, maxLength  
Added operation descriptions  
Added properties description  
Example now matching pattern  
Binary values use substituted with url for png  
OperationId Checked  
Fixed likes' design  
Added Security scheme  
minItems and maxItems added to array use  
snake_case notation  
For security reasons and according to IBM validation standards all methods that should return an array will return a JSON containing an object type resource with only "list" property, built as an array of objects.  
setMyUserName now uses put method instead of patch method  
followUser now take identifier instead of username as schema  
Every parameter/schema/property has now its own description  
Change from /user/{uid}/following/stream to /user/{uid}/stream  
Change from followed_uid to following_uid  
getUserProfile now operate in /user/{uid} and profile schema has new property "identifier"  

