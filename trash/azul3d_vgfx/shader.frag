#version 120

varying vec3 vInterp;
//pow(u,3) - u*v = u(pow(u,2) - v)

float quadraticBezier(vec2 p) {
	// Gradients
	vec2 px = dFdx(p);
	vec2 py = dFdy(p);

	// Chain rule
	float fx = (2*p.x)*px.x - px.y;
	float fy = (2*p.x)*py.x - py.y;

	// Signed distance
	float sd = (p.x*p.x - p.y)/sqrt(fx*fx + fy*fy);

	// Linear alpha
	return 0.5 - sd;
}

void main()
{
	float alpha = quadraticBezier(vInterp.xz);
	gl_FragColor = vec4(0,1,0,1);

	if(alpha < 1) {
		return;
	} else if(alpha > 0) {
		discard;
		return;
	} else {
		gl_FragColor.a = alpha;
	}

/*
	if(alpha > 1) {
		return;
	} else if(alpha < 0) {
		discard;
		return;
	} else {
		gl_FragColor.a = alpha;
	}
*/
}

