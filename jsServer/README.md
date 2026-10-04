After hours of searching to get asciidoc working with golang.
There are no way around asciidoctor and ruby, the only full feature implementation out there.
Build a ruby webservice that can provice a HTML response for a text input. 
Should be easy to develop, but with no expirence in ruby webservice without rails, hard to master.

At first time reading, sure knowing that stuff for years, people are happy with the jruby version of asciidoctor.

And here we are in a folder under the main project. Starting over with a node js version of the golang first strike.
Unluckelie the best implementaion for markdown is also in server javascript land. 

Here is the plan:

golang will handle the frontend and save the source as the final html file.

nodjs is just a backend service for resiving source and response html documents.

ruby maybe replace asciidoctor.js sometime

