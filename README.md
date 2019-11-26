## Ray the Tracer

In graduate school I took a 3D rendering class and wrote a ray tracing
renderer in C++ using the Qt libraries for drawing.  The class was
great and I learned a lot, but sadly I did a poor job commenting my
code, and was a little overzealous with my operator overloading.  I
also have not kept up on my C++ so reading through it now is like
reading a foreign language.  I have no idea where my algorithms came
from and no idea how to add new objects to the environment or really
manipulate anything.

So, 10+ years later I'm giving it another try only this time I'm going
to try writing it in Go so I can take advantage of the parallel
processing, memory management and profiling tools.  I'm also specifically 
not taking shortcuts for the sake of clarity.

The ray tracer in the [Soda Water](cmd/sodaWater) folder 
is my independent attempt to build a ray tracer.  I got part way through 
this project and realized that I was spending too much time re-deriving 
all of the math equations and decided to instead use a reference book. 

The ray tracer in the [Weekend](cmd/weekend) folder is my attempt to build a 
ray tracer from the book
[Ray Tracing in a Weekend](https://www.realtimerendering.com/raytracing/Ray%20Tracing%20in%20a%20Weekend.pdf)
