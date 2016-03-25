#version 120

varying vec3 vBC;

void main()
{
	vec3 v = vBC;

	float yg = 1 - abs(v.z - v.y); // Vertical Gradient
	float xg = v.x;                // Horizontal Gradient
	float t  = .5 - (v.y - v.x)/2; // T value

	float a1 = .5 - (v.y - v.x)/2;
	float a2 = .5 - (v.z - v.x)/2;
	float bt = mix(a1, xg, 1-t);
	bt = a1;

	gl_FragColor = vec4(bt,bt,bt,1);
	if(bt > 0.5) {
		discard;
	}
}

