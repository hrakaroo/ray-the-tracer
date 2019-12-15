TODOs
=====

* ~~Colored transparent objects~~
* ~~Block rotation~~
* Lights and shadows
* Caustics
* Camera up -> tilt

Ponderings
==========

Should New* functions return objects or pointers?
-------------------------------------------------

In general I'm not a huge fan of Go's approach here.  If we don't use New constructors
then adding additional fields to structs is asking for errors as there is no way to
enforce/check that the new field is getting set everywhere the struct is getting created.
So instead we can create New constructors for every struct, but at that point why not
just support constructors on objects.  This seems to be another place that Go started with
an idealized approach which then broke down as soon as people actually started using it.

Regardless, the real question is if the New constructor should return an object or
a reference.  I like consistency so initially I tried to always return an object, 
but that runs quickly into an issue where New is used as an argument for a 
function that wants a reference as you _can't_ take the reference without first
assigning it to an intermediate variable.  (I'm also not a fan of this.)  
But swinging entirely the other direction doesn't seem correct either as it 
complicates the functions on Vec3.  

So one solution might be to create two methods, one called New* and the other
New*Ref.  This would at least be consistent, but ugly.  Or, we could set up
a rule that "primitives" return objects but complex structures return references.
So Vec3 and Ray both return objects, but Sphere and World are references (?)

I'm still trying to figure out the best approach here as I don't like making it up
on the fly as I will invariably create the right usage for what I'm working on 
_at that moment_ but will end up causing chaos later.


What to return?
---------------

I'm not a huge fan of having multiple ways of sending data back to the caller and 
with Go's ability to return multiple values I lean heavily against returning data
by passing in references to objects.  That said, I'm also not hugely keen on creating
a datatype just to ignore it because either error is not nil or some boolean is false.
This often pushes me to always return pointers so I can easily return 
`false, nil`, as opposed to `false, NewStruct(bogusValues)`.  That said, I'm 
returning an empty Vec3 now to represent dark so ... perhaps it's not worth 
worrying about.