# WasaPhoto
Share your flash!

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
---DONE  
Common path parameters fixed  
---DONE  
Defined pattern, minLength, maxLength  
---DONE  
Added operation descriptions  
---DONE  
Added properties description  
---DONE  
Example now matching pattern  
---DONE  
Binary values use substituted with url for png  
---DONE  
OperationId Checked  
---DONE  
Fixed likes' design  
---UNDONE (POSSIBLE BRENCH TO REDESIGN)  
Added Security scheme  
---DONE  
minItems and maxItems added to array use  
---DONE  
snake_case notation  
---DONE  
  
-For security reasons and according to IBM validation standards all methods that should return an array will return a JSON containing an object type resource with only
"list" property, built as an array of objects.  
-setMyUserName now uses put method instead of patch method  
-followUser now take identifier instead of username as schema  
-Every parameter/schema/property has now its own description  
-Change from /user/{uid}/following/stream to /user/{uid}/stream  
-Change from followed_uid to following_uid  
-getUserProfile now operate in /user/{uid} and profile schema has new property "identifier"  


