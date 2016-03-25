#version 120

attribute vec3 Vertex;
attribute vec3 Interp;

uniform mat4 MVP;

varying vec3 vInterp;

void main()
{
	vInterp = Interp;
	gl_Position = MVP * vec4(Vertex, 1.0);
}

