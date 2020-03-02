# go3D: A software-based real time renderer in Go

![A rendered frame](output.png)

## What?

go3D is an on-CPU rasterizing 3D rendering engine written in the Go language. As an education project, it does so without the use of standard graphics frameworks like OpenGL, DirectX, or Vulkan. Though this project does relies on several external libraries, dependencies are only used to provide basic functionality, with function such as 2D line drawing and polygon filling, matrix multiplication, and file I/O handled by the main project. It uses performance-optimized features, including threads, channels, and slices.

### Features

go3D supports several common features of basic 3D programs:

- Reading from .obj files
- Different perspectives and viewing fulcrums, including perspective and orthographic projections
- Wireframe and solid shading
- Backface culling
- Z-buffering
- Camera movement
- Basic animations, including spinning and moving

## Why?

This project was originally created as an independent research project for a high-school research to study computer graphics, performance optimization, memory management, and the Go language. Some minor modifications have since been performed in an effort to better understand the mathematical foundations of computer graphics and perform minor cleanup.

## Building and running

If you really, really want to run this, you can! That said, since this is an educational project, is is not meant to be distributed. You may find it more informative to just poke around the code instead.

1. Install [gonum](https://github.com/gonum/gonum) and [Fyne](https://github.com/fyne-io/fyne). 
2. Configure and load the models you'd like to view (a default is provided.)
 - Configuration is currently done by simply modifying the objects in the example `main.go` file. Commented-out examples are (loosely) provided.
 - Models must be properly triangulated. It is recommended to run all models through a program such as Blender prior to loading them in.
3. 	Set the transformations you'd like to apply.
 - Loose examples are again provided. 
4. Configure the view.
 - You can change the viewing window size, the scaling, the camera position, and the projection matrix. You can also leave them as default.
5. Run!

## Other documentation

If you'd like to read up on the mathematics behind the rendering, **a paper is provided in `doc/` detailing the applications of linear algebra to computer graphics.** Note that since this paper assumes familiarity with linear algebra.


## Disclaimers

This code is provided as-is with no guarantee of functionality. Development is considered stopped as the academic purpose of this project has been fulfilled for the time being. Don't worry about the Git log.

Unfortunately, due to the academic nature of this project, restrictions on the distribution of this code must currently be imposed due to academic integrity policies. Please do not redistribute this code without my express permission. That said, if you are interested, please contact me via this account and I'll be happy to add you to the repository. Hopefully this can change in the future.