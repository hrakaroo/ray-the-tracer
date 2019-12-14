More Images
===========

This is just a set of additional images I've generated as I continue to develop the ray
tracer.  I have a set of additional things I want to explore before I'm done with this 
project and move onto my next.

---

This is the one image I created in my sodaWater ray tracer before I abandoned that
effort and started working with the Ray Tracer in a Weekend book.  It's a block, but
my lighting was all messed up.

![Image](initial.png)

---

![Image](blocks01.png)

However, I was able to take the code for rendering blocks from the other ray tracer
and mostly copy it directly into my weekend project.  Thankfully I still remembered 
the basic algorithm for how to do this from when I took the 3d graphics class in 
graduate school.  I really like how the blocks look in this scene.

---

All of the other images use a huge sphere to represent the ground, but for this I 
thought I would replace it with a real block and I made it metallic as well.  
I really like this picture, but I do miss the shadows which don't render as well
 on metal.

![Image](blocks02.png)

---

Just a big image with blocks and using a big block as the ground.  I don't love this
image, but I am using it as a test to see if running multiple goroutines can speed
up the render.  This image was generated with 50 goroutines and took 455m to finish.

![Image](blocks03.png)

---

As I was explaining to someone how rendering a block works, it occurred to me that
it should be easy to detect the edges of a block and outline them in black.  This was
my first attempt to test it out and it worked really well.  I don't think I'll keep
it on as a default, but it was pretty cool.

![Image](blocks04.png)


