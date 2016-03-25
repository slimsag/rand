#version 120

attribute vec3 Vertex;
attribute vec3 Bary;

uniform mat4 MVP;

varying vec3 vBC;

void main()
{
	vBC = Bary;
	gl_Position = MVP * vec4(Vertex, 1.0);
}

