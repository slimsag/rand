package main

import (
	"azul3d.org/v1/gfx"
)

// wireVert is GLSL vertex shader source code for a basic wireframe shader.
var wireVert = []byte(`
#version 120

attribute vec3 Vertex;
attribute vec3 Bary;
varying vec3 vBC;
uniform mat4 MVP;

void main()
{
	vBC = Bary;
	gl_Position = MVP * vec4(Vertex, 1.0);
}
`)

// wireFrag is GLSL fragment shader source code for a basic wireframe shader.
var wireFrag = []byte(`
#version 120
#extension GL_OES_standard_derivatives: enable

#define FRONT_LINE_WIDTH 2.0
#define BACK_LINE_WIDTH 1.0
#define FRONT_COLOR vec3(0.0, 0.5, 0.0)
#define BACK_COLOR vec3(1.0, 0.0, 0.0)

varying vec3 vBC;
uniform bool BinaryAlpha;

float edgeFactor(float lineWidth) {
	vec3 d = fwidth(vBC);
	vec3 a3 = smoothstep(vec3(0.0), d*lineWidth, vBC);
	return min(min(a3.x, a3.y), a3.z);
}

void main() {
	if(gl_FrontFacing){
		gl_FragColor = vec4(FRONT_COLOR, (1.0-edgeFactor(FRONT_LINE_WIDTH)));
	} else{
		gl_FragColor = vec4(BACK_COLOR, (1.0-edgeFactor(BACK_LINE_WIDTH)));
	}
	if(BinaryAlpha && gl_FragColor.a < 0.5) {
		discard;
	}
}
`)

var Wireframe *gfx.Shader

func init() {
	Wireframe = gfx.NewShader("Wireframe Shader")
	Wireframe.GLSLVert = wireVert
	Wireframe.GLSLFrag = wireFrag
}
