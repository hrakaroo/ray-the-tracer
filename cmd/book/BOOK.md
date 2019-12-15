Ray Tracing in a Weekend
========================

The book I'm using as my reference/tutorial is [Ray Tracing in a Weekend](https://www.realtimerendering.com/raytracing/Ray%20Tracing%20in%20a%20Weekend.pdf).
Although the author provides it for free I would encourage buying a digital copy from Amazon for 
$2.99.  It's well worth the cost.

I'm specifically not trying to just copy/translate the code from the C++ in the book into Go. 
Instead I'm trying to understand each section and rewrite it in a way that both makes sense
to me and is more inline with the Go style. 

For example, I'm going to try to take advantage of the multiple return types that Go
supports and specifically not use operator overloading (partly because Go doesn't support
them, but also I find them to be nearly impossible to read later.)

That said, I do have a couple of issues/questions with the book.  What follows is an unordered
list.


1. Unused variable
------------------
In Chapter 6, Antialiasing, the Main is changed to
```c
  ray r = cam.get_ray(u, v);
  vec3 p = r.point_at_parameter(2.0);
  col += color(r, world);
```

Point `p` doesn't appear to be used anywhere.  I'm guessing this is just a mistake.


2. Eliminated a bunch of redundant 2's
---------------------------------------
In Chapter 5, Surface normals and multiple objects, the quadratic equation is used
to determine the value of x as
```c
float temp = (-b - sqrt(b*b-a*c))/a;
```

I need to actually derive this to verify to myself that this is actually equal to
```c
float temp = (-b - sqrt(b*b - 4*a*c))/2*a
```

It might actually be correct, especially as we are often dealing with unit vectors, but
it smells a little funny to me.


3. Unit cube
------------
In Chapter 7, Diffuse Materials it says to "pick a random point in the unit cube where x,y, 
and z all range from -1 to +1".  If I'm not mistaken, that is not the unit cube.  A 
unit cube would be from -0.5 to + 0.5 for a total length of 1 unit on each side.


4. Operator overloading
-----------------------
This is a general comment that I'm not a fan of operator overloading.  If you _live_ the
code then I guess, but if you are like most people who are going to write this and then 
not revisit it for months, the operator overloading becomes almost unreadable.  And worse, 
most IDEs don't support control clicking into an operator like you can with a basic method
so even navigating it is painful.  Not a fan.