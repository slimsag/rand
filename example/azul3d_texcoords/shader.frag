#version 120

varying vec2 tc0;

uniform sampler2D Texture0;
uniform bool BinaryAlpha;
uniform bool Enabled;

float side(vec2 lp1, vec2 lp2, vec2 p) {
	return (p.x - lp1.x) * (lp2.y - lp1.y) - (p.y - lp1.y) * (lp2.x - lp1.x);
}

vec2 pixelToCoord(vec2 invImageSize, vec2 coord) {
	return invImageSize * coord;
}

vec2 coordToPixel(vec2 invImageSize, vec2 coord) {
	return coord / invImageSize;
}

vec4 sharpSample(sampler2D tex, vec2 invTexSize, vec2 texCoord, vec4 passThrough) {
	vec2 texCoordPixel = coordToPixel(invTexSize, texCoord);
	vec2 cp = floor(texCoordPixel) + 0.5;

	vec4 center        = texture2D(tex, texCoord);
	vec4 left          = texture2D(tex, pixelToCoord(invTexSize, vec2(cp.x-1, cp.y)));
	vec4 right         = texture2D(tex, pixelToCoord(invTexSize, vec2(cp.x+1, cp.y)));
	vec4 bottom        = texture2D(tex, pixelToCoord(invTexSize, vec2(cp.x,   cp.y+1)));
	vec4 top           = texture2D(tex, pixelToCoord(invTexSize, vec2(cp.x,   cp.y-1)));

	if(bottom == right && left == top && top != right) {
		vec2 linePoints = vec2(cp.x - 0.5, cp.y + 0.5);
		if(side(cp, linePoints, texCoordPixel) > 0) {
			return bottom;
		} else {
			return passThrough;
		}
	}
	return passThrough;
}

float rand(vec2 co){
    return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

float cum(vec4 v) {
	return (v.x + v.y + v.z + v.w) / 4.0;
}

void main()
{
	gl_FragColor = texture2D(Texture0, tc0);
	if(BinaryAlpha && gl_FragColor.a < 0.5) {
		discard;
	}

	if(!Enabled) {
		return;
	}

	vec4 p = texture2D(Texture0, tc0);

	int samples = 0;
	float sampleSize = 1/128.0;
	float selection = .6;

	vec4 sample = vec4(1, 1, 1, 1);
	float avgSample = 0.0;
	for(int i=0; i < samples; i++) {
		float r1 = rand(tc0) - 0.5;
		float r2 = rand(tc0 + vec2(10, 10)) - 0.5;
		vec2 offset = vec2(sampleSize * r1, sampleSize * r2);
		vec4 s = texture2D(Texture0, clamp(tc0 + offset, 0, 1));
		avgSample += cum(s);
	}
	if(samples > 0) {
		avgSample = avgSample / samples;
		selection = (selection + avgSample) / 2.0;
	}

	if(any(lessThan(p, vec4(selection, selection, selection, 1.0)))) {
		gl_FragColor = vec4(0, 0, 0, 1);
	} else {
		discard;
	}

	//vec2 invTexSize = 1.0 / vec2(32.0, 32.0);
	//gl_FragColor = sharpSample(Texture0, invTexSize, tc0, gl_FragColor);

	return;

	/*
	if(topRight != gl_FragColor && right != gl_FragColor) {
		//float dist = distance(coordToPixel(invImageSize, tc0), coordToPixel(invImageSize, topRightCoord));
		float dist = distance(coordToPixel(invImageSize, tc0), vec2(0, 0));
		float margin = 16;
		if(dist < margin) {
			gl_FragColor = vec4(0.0, 1.0, 0.0, 1.0);
		} else {
			gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
		}
	}
	*/
}

