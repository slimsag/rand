package ai

/*
#include "assimp/config.h"
#include "assimp/postprocess.h"
*/
import "C"

// Enumerates components of the Scene and Mesh data structures
// that can be excluded from the import using the #aiPrpcess_RemoveComponent step.
//
// See the documentation to #aiProcess_RemoveComponent for more details.
type Component int

// Remove a specific color channel 'n'
func COLORSn(n uint) Component {
	//aiComponent_COLORSn(n) (1u << (n+20u))
	return Component(1 << (n + 20))
}

// Remove a specific UV channel 'n'
func TEXCOORDSn(n uint) Component {
	//aiComponent_TEXCOORDSn(n) (1u << (n+25u))
	return Component(1 << (n + 25))
}

const (
	// Normal vectors
	NORMALS = C.aiComponent_NORMALS

	// Tangents and bitangents go always together ...
	TANGENTS_AND_BITANGENTS = C.aiComponent_TANGENTS_AND_BITANGENTS

	// ALL color sets
	// Use aiComponent_COLORn(N) to specify the N'th set
	COLORS = C.aiComponent_COLORS

	// ALL texture UV sets
	// aiComponent_TEXCOORDn(N) to specify the N'th set
	TEXCOORDS = C.aiComponent_TEXCOORDS

	// Removes all bone weights from all meshes.
	// The scenegraph nodes corresponding to the bones are NOT removed.
	// use the #aiProcess_OptimizeGraph step to do this
	BONEWEIGHTS = C.aiComponent_BONEWEIGHTS

	// Removes all node animations (aiScene::mAnimations).
	// The corresponding scenegraph nodes are NOT removed.
	// use the #aiProcess_OptimizeGraph step to do this
	ANIMATIONS = C.aiComponent_ANIMATIONS

	// Removes all embedded textures (aiScene::mTextures)
	TEXTURES = C.aiComponent_TEXTURES

	// Removes all light sources (aiScene::mLights).
	// The corresponding scenegraph nodes are NOT removed.
	// use the #aiProcess_OptimizeGraph step to do this
	LIGHTS = C.aiComponent_LIGHTS

	// Removes all cameras (aiScene::mCameras).
	// The corresponding scenegraph nodes are NOT removed.
	// use the #aiProcess_OptimizeGraph step to do this
	CAMERAS = C.aiComponent_CAMERAS

	// Removes all meshes (aiScene::mMeshes).
	MESHES = C.aiComponent_MESHES

	// Removes all materials. One default material will
	// be generated, so aiScene::mNumMaterials will be 1.
	MATERIALS = C.aiComponent_MATERIALS
)

type Prop string

const (
	// Enables time measurements.
	//
	// If enabled, measures the time needed for each part of the loading
	// process (i.e. IO time, importing, postprocessing, ..) and dumps
	// these timings to the DefaultLogger. See the  perf Performance
	// Page for more information on this topic.
	//
	// Property type: bool. Default value: false.
	GlobMeasureTime Prop = "GLOB_MEASURE_TIME"

	// Global setting to disable generation of skeleton dummy meshes
	//
	// Skeleton dummy meshes are generated as a visualization aid in cases which
	// the input data contains no geometry, but only animation data.
	// Property data type: bool. Default value: false
	ImportNoSkeletonMeshes Prop = "IMPORT_NO_SKELETON_MESHES"

	// Maximum bone count per mesh for the SplitbyBoneCount step.
	//
	// Meshes are split until the maximum number of bones is reached. The default
	// value is AI_SBBC_DEFAULT_MAX_BONES, which may be altered at
	// compile-time.
	// Property data type: integer.
	SBBCMaxBones Prop = "PP_SBBC_MAX_BONES"

	// default limit for bone count
	SBBCDefaultMaxBones = 60

	// Specifies the maximum angle that may be between two vertex tangents
	//   that their tangents and bi-tangents are smoothed.
	//
	// This applies to the CalcTangentSpace-Step. The angle is specified
	// in degrees. The maximum value is 175.
	// Property type: float. Default value: 45 degrees
	CTMaxSmoothingAngle Prop = "PP_CT_MAX_SMOOTHING_ANGLE"

	// Source UV channel for tangent space computation.
	//
	// The specified channel must exist or an error will be raised.
	// Property type: integer. Default value: 0
	CTTextureChannelIndex Prop = "PP_CT_TEXTURE_CHANNEL_INDEX"

	// Specifies the maximum angle that may be between two face normals
	//    at the same vertex position that their are smoothed together.
	//
	// Sometimes referred to as 'crease angle'.
	// This applies to the GenSmoothNormals-Step. The angle is specified
	// in degrees, so 180 is PI. The default value is 175 degrees (all vertex
	// normals are smoothed). The maximum value is 175, too. Property type: float.
	// Warning: setting this option may cause a severe loss of performance. The
	// performance is unaffected if the #AI_CONFIG_FAVOUR_SPEED flag is set but
	// the output quality may be reduced.
	GSNMaxSmoothingAngle Prop = "PP_GSN_MAX_SMOOTHING_ANGLE"

	// Sets the colormap (= palette) to be used to decode embedded
	//   textures in MDL (Quake or 3DGS) files.
	//
	// This must be a valid path to a file. The file is 768 (256// 3) bytes
	// large and contains RGB triplets for each of the 256 palette entries.
	// The default value is colormap.lmp. If the file is not found,
	// a default palette (from Quake 1) is used.
	// Property type: string.
	ImportMDLColormap Prop = "IMPORT_MDL_COLORMAP"

	// Configures the #aiProcess_RemoveRedundantMaterials step to
	// keep materials matching a name in a given list.
	//
	// This is a list of 1 to n strings, ' ' serves as delimiter character.
	// Identifiers containing whitespaces must be enclosed in// single//
	// quotation marks. For example:<tt>
	// "keep-me and_me_to anotherMaterialToBeKept \'name with whitespace\'"</tt>.
	// If a material matches on of these names, it will not be modified or
	// removed by the postprocessing step nor will other materials be replaced
	// by a reference to it. <br>

	// This option might be useful if you are using some magic material names
	// to pass additional semantics through the content pipeline. This ensures
	// they won't be optimized away, but a general optimization is still
	// performed for materials not contained in the list.

	// Property type: String. Default value: n/a
	// @note Linefeeds, tabs or carriage returns are treated as whitespace.
	// Material names are case sensitive.
	RRMExcludeList Prop = "PP_RRM_EXCLUDE_LIST"

	// Configures the #aiProcess_PretransformVertices step to
	// keep the scene hierarchy. Meshes are moved to worldspace, but
	// no optimization is performed (read: meshes with equal materials are not
	// joined. The total number of meshes won't change).
	//
	// This option could be of use for you if the scene hierarchy contains
	// important additional information which you intend to parse.
	// For rendering, you can still render all meshes in the scene without
	// any transformations.
	// Property type: bool. Default value: false.
	PTVKeepHierarchy Prop = "PP_PTV_KEEP_HIERARCHY"

	// Configures the #aiProcess_PretransformVertices step to normalize
	// all vertex components into the [-1,1] range. That is, a bounding box
	// for the whole scene is computed, the maximum component is taken and all
	// meshes are scaled appropriately (uniformly of course!).
	// This might be useful if you don't know the spatial dimension of the input
	// data
	PTVNormalize Prop = "PP_PTV_NORMALIZE"

	// Configures the #aiProcess_FindDegenerates step to
	// remove degenerated primitives from the import - immediately.
	//
	// The default behaviour converts degenerated triangles to lines and
	// degenerated lines to points. See the documentation to the
	// #aiProcess_FindDegenerates step for a detailed example of the various ways
	// to get rid of these lines and points if you don't want them.
	// Property type: bool. Default value: false.
	FDRemove Prop = "PP_FD_REMOVE"

	// Configures the #aiProcess_OptimizeGraph step to preserve nodes
	// matching a name in a given list.
	//
	// This is a list of 1 to n strings, ' ' serves as delimiter character.
	// Identifiers containing whitespaces must be enclosed in// single//
	// quotation marks. For example:<tt>
	// "keep-me and_me_to anotherNodeToBeKept \'name with whitespace\'"</tt>.
	// If a node matches on of these names, it will not be modified or
	// removed by the postprocessing step.<br>
	// This option might be useful if you are using some magic node names
	// to pass additional semantics through the content pipeline. This ensures
	// they won't be optimized away, but a general optimization is still
	// performed for nodes not contained in the list.
	// Property type: String. Default value: n/a
	// @note Linefeeds, tabs or carriage returns are treated as whitespace.
	// Node names are case sensitive.
	OGExcludeList Prop = "PP_OG_EXCLUDE_LIST"

	// Set the maximum number of triangles in a mesh.
	//
	// This is used by the "SplitLargeMeshes" PostProcess-Step to determine
	// whether a mesh must be split or not.
	// @note The default value is AI_SLM_DEFAULT_MAX_TRIANGLES
	// Property type: integer.
	SLMTriangleLimit Prop = "PP_SLM_TRIANGLE_LIMIT"

	// default value for SLMTriangleLimit
	SLMDefaultMaxTriangle = 1000000

	// Set the maximum number of vertices in a mesh.
	//
	// This is used by the "SplitLargeMeshes" PostProcess-Step to determine
	// whether a mesh must be split or not.
	// @note The default value is AI_SLM_DEFAULT_MAX_VERTICES
	// Property type: integer.
	SLMVertexLimit Prop = "PP_SLM_VERTEX_LIMIT"

	// default value for SLMVertexLimit
	SLMDefaultMaxVertices = 1000000

	// Set the maximum number of bones affecting a single vertex
	//
	// This is used by the #aiProcess_LimitBoneWeights PostProcess-Step.
	// @note The default value is AI_LBW_MAX_WEIGHTS
	// Property type: integer.
	LBWMaxWeights Prop = "PP_LBW_MAX_WEIGHTS"

	// default value for LBWMaxWeights
	LBWDefaultMaxWeights = 0x4

	// Lower the deboning threshold in order to remove more bones.
	//
	// This is used by the #aiProcess_Debone PostProcess-Step.
	// @note The default value is AI_DEBONE_THRESHOLD
	// Property type: float.
	DBThreshold Prop = "PP_DB_THRESHOLD"

	// default value for DBThreshold
	DBDefaultThreshold float32 = 1.0

	// Require all bones qualify for deboning before removing any
	//
	// This is used by the #aiProcess_Debone PostProcess-Step.
	// @note The default value is 0
	// Property type: bool.
	DBAllOrNone Prop = "PP_DB_ALL_OR_NONE"

	// Default value for the ICLPtcacheSize
	ICLDefaultPtcacheSize = 12

	// Set the size of the post-transform vertex cache to optimize the
	// vertices for. This configures the #aiProcess_ImproveCacheLocality step.
	//
	// The size is given in vertices. Of course you can't know how the vertex
	// format will exactly look like after the import returns, but you can still
	// guess what your meshes will probably have.
	// @note The default value is #PP_ICL_PTCACHE_SIZE. That results in slight
	// performance improvements for most nVidia/AMD cards since 2002.
	// Property type: integer.
	ICLPtcacheSize Prop = "PP_ICL_PTCACHE_SIZE"

	// Input parameter to the #aiProcess_RemoveComponent step:
	// Specifies the parts of the data structure to be removed.
	//
	// See the documentation to this step for further details. The property
	// is expected to be an integer, a bitwise combination of the
	// #aiComponent flags defined above in this header. The default
	// value is 0. Important: if no valid mesh is remaining after the
	// step has been executed (e.g you thought it was funny to specify ALL
	// of the flags defined above) the import FAILS. Mainly because there is
	// no data to work on anymore ...
	RVCFlags Prop = "PP_RVC_FLAGS"

	// Input parameter to the #aiProcess_SortByPType step:
	// Specifies which primitive types are removed by the step.
	//
	// This is a bitwise combination of the aiPrimitiveType flags.
	// Specifying all of them is illegal, of course. A typical use would
	// be to exclude all line and point meshes from the import. This
	// is an integer property, its default value is 0.
	SBPRemove Prop = "PP_SBP_REMOVE"

	// Input parameter to the #aiProcess_FindInvalidData step:
	// Specifies the floating-point accuracy for animation values. The step
	// checks for animation tracks where all frame values are absolutely equal
	// and removes them. This tweakable controls the epsilon for floating-point
	// comparisons - two keys are considered equal if the invariant
	// abs(n0-n1)>epsilon holds true for all vector respectively quaternion
	// components. The default value is 0.f - comparisons are exact then.
	FIDAnimAccuracy Prop = "PP_FID_ANIM_ACCURACY"

	// TransformUVCoords evaluates UV scalings
	UVTRAFO_SCALING = C.AI_UVTRAFO_SCALING

	// TransformUVCoords evaluates UV rotations
	UVTRAFO_ROTATION = C.AI_UVTRAFO_ROTATION

	// TransformUVCoords evaluates UV translation
	UVTRAFO_TRANSLATION = C.AI_UVTRAFO_TRANSLATION

	// Everything baked together -> default value
	UVTRAFO_ALL = C.AI_UVTRAFO_ALL

	// Input parameter to the #aiProcess_TransformUVCoords step:
	// Specifies which UV transformations are evaluated.
	//
	// This is a bitwise combination of the AI_UVTRAFO_XXX flags (integer
	// property, of course). By default all transformations are enabled
	// (AI_UVTRAFO_ALL).
	TUVEvaluate Prop = "PP_TUV_EVALUATE"

	// A hint to assimp to favour speed against import quality.
	//
	// Enabling this option may result in faster loading, but it needn't.
	// It represents just a hint to loaders and post-processing steps to use
	// faster code paths, if possible.
	// This property is expected to be an integer, != 0 stands for true.
	// The default value is 0.
	FavourSpeed Prop = "FAVOUR_SPEED"

	// ###########################################################################
	// IMPORTER SETTINGS
	// Various stuff to fine-tune the behaviour of specific importer plugins.
	// ###########################################################################

	// Set whether the fbx importer will merge all geometry layers present
	// in the source file or take only the first.
	//
	// The default value is true (1)
	// Property type: bool
	FBXReadAllGeometryLayers Prop = "IMPORT_FBX_READ_ALL_GEOMETRY_LAYERS"

	// Set whether the fbx importer will read all materials present in the
	// source file or take only the referenced materials.
	//
	// This is void unless IMPORT_FBX_READ_MATERIALS=1.
	//
	// The default value is false (0)
	// Property type: bool
	FBXReadAllMaterials Prop = "IMPORT_FBX_READ_ALL_MATERIALS"

	// Set whether the fbx importer will read materials.
	//
	// The default value is true (1)
	// Property type: bool
	FBXReadMaterials Prop = "IMPORT_FBX_READ_MATERIALS"

	// Set whether the fbx importer will read cameras.
	//
	// The default value is true (1)
	// Property type: bool
	FBXReadCameras Prop = "IMPORT_FBX_READ_CAMERAS"

	// Set whether the fbx importer will read light sources.
	//
	// The default value is true (1)
	// Property type: bool
	FBXReadLights Prop = "IMPORT_FBX_READ_LIGHTS"

	// Set whether the fbx importer will read animations.
	//
	// The default value is true (1)
	// Property type: bool
	FBXReadAnimations Prop = "IMPORT_FBX_READ_ANIMATIONS"

	// Set whether the fbx importer will act in strict mode in which only
	// FBX 2013 is supported and any other sub formats are rejected. FBX 2013
	// is the primary target for the importer, so this format is best
	// supported and well-tested.
	//
	// The default value is false (0)
	// Property type: bool
	FBXStrictMode Prop = "IMPORT_FBX_STRICT_MODE"

	// Set whether the fbx importer will preserve pivot points for
	// transformations (as extra nodes). If set to false, pivots and offsets
	// will be evaluated whenever possible.
	//
	// The default value is true (1)
	// Property type: bool
	FBXPreservePivots Prop = "IMPORT_FBX_PRESERVE_PIVOTS"

	// Specifies whether the importer will drop empty animation curves or
	// animation curves which match the bind pose transformation over their
	// entire defined range.
	//
	// The default value is true (1)
	// Property type: bool
	FBXOptimizeEmptyAnimationCurves Prop = "IMPORT_FBX_OPTIMIZE_EMPTY_ANIMATION_CURVES"

	// Set the vertex animation keyframe to be imported
	//
	// ASSIMP does not support vertex keyframes (only bone animation is supported).
	// The library reads only one frame of models with vertex animations.
	// By default this is the first frame.
	// \note The default value is 0. This option applies to all importers.
	// However, it is also possible to override the global setting
	// for a specific loader. You can use the AI_CONFIG_IMPORT_XXX_KEYFRAME
	// options (where XXX is a placeholder for the file format for which you
	// want to override the global setting).
	// Property type: integer.
	GlobalKeyframe Prop = "IMPORT_GLOBAL_KEYFRAME"

	MD3Keyframe    Prop = "IMPORT_MD3_KEYFRAME"
	MD2Keyframe    Prop = "IMPORT_MD2_KEYFRAME"
	MDLKeyframe    Prop = "IMPORT_MDL_KEYFRAME"
	MDCKeyframe    Prop = "IMPORT_MDC_KEYFRAME"
	SMDKeyframe    Prop = "IMPORT_SMD_KEYFRAME"
	UnrealKeyframe Prop = "IMPORT_UNREAL_KEYFRAME"

	// Configures the AC loader to collect all surfaces which have the
	// "Backface cull" flag set in separate meshes.
	//
	// Property type: bool. Default value: true.
	ACSeparateBFCull Prop = "IMPORT_AC_SEPARATE_BFCULL"

	// Configures whether the AC loader evaluates subdivision surfaces (
	// indicated by the presence of the 'subdiv' attribute in the file). By
	// default, Assimp performs the subdivision using the standard
	// Catmull-Clark algorithm
	//
	// Property type: bool. Default value: true.
	ACEvalSubdivision Prop = "IMPORT_AC_EVAL_SUBDIVISION"

	// Configures the UNREAL 3D loader to separate faces with different
	// surface flags (e.g. two-sided vs. single-sided).
	//
	// Property type: bool. Default value: true.
	UnrealHandleFlags Prop = "UNREAL_HANDLE_FLAGS"

	// Configures the terragen import plugin to compute uv's for
	// terrains, if not given. Furthermore a default texture is assigned.
	//
	// UV coordinates for terrains are so simple to compute that you'll usually
	// want to compute them on your own, if you need them. This option is intended
	// for model viewers which want to offer an easy way to apply textures to
	// terrains.
	//
	// Property type: bool. Default value: false.
	TERMakeUVS Prop = "IMPORT_TER_MAKE_UVS"

	// Configures the ASE loader to always reconstruct normal vectors
	// 	basing on the smoothing groups loaded from the file.
	//
	// Some ASE files have carry invalid normals, other don't.
	//
	// Property type: bool. Default value: true.
	ASEReconstructNormals Prop = "IMPORT_ASE_RECONSTRUCT_NORMALS"

	// Configures the M3D loader to detect and process multi-part
	// Quake player models.
	//
	// These models usually consist of 3 files, lower.md3, upper.md3 and
	// head.md3. If this property is set to true, Assimp will try to load and
	// combine all three files if one of them is loaded.
	//
	// Property type: bool. Default value: true.
	MD3HandleMultipart Prop = "IMPORT_MD3_HANDLE_MULTIPART"

	// Tells the MD3 loader which skin files to load.
	//
	// When loading MD3 files, Assimp checks whether a file
	// <md3_file_name>_<skin_name>.skin is existing. These files are used by
	// Quake III to be able to assign different skins (e.g. red and blue team)
	// to models. 'default', 'red', 'blue' are typical skin names.
	//
	// Property type: String. Default value: "default".
	MD3SkinName Prop = "IMPORT_MD3_SKIN_NAME"

	// Specify the Quake 3 shader file to be used for a particular
	// MD3 file. This can also be a search path.
	//
	// By default Assimp's behaviour is as follows: If a MD3 file
	// <tt><any_path>/models/<any_q3_subdir>/<model_name>/<file_name>.md3</tt> is
	// loaded, the library tries to locate the corresponding shader file in
	// <tt><any_path>/scripts/<model_name>.shader</tt>. This property overrides this
	// behaviour. It can either specify a full path to the shader to be loaded
	// or alternatively the path (relative or absolute) to the directory where
	// the shaders for all MD3s to be loaded reside. Assimp attempts to open
	// <tt><dir>/<model_name>.shader</tt> first, <tt><dir>/<file_name>.shader</tt>
	// is the fallback file. Note that <dir> should have a terminal (back)slash.
	//
	// Property type: String. Default value: n/a.
	MD3ShaderSrc Prop = "IMPORT_MD3_SHADER_SRC"

	// Configures the LWO loader to load just one layer from the model.
	//
	// LWO files consist of layers and in some cases it could be useful to load
	// only one of them. This property can be either a string - which specifies
	// the name of the layer - or an integer - the index of the layer. If the
	// property is not set the whole LWO model is loaded. Loading fails if the
	// requested layer is not available. The layer index is zero-based and the
	// layer name may not be empty.<br>
	//
	// Property type: Integer. Default value: all layers are loaded.
	LWOOneLayerOnly Prop = "IMPORT_LWO_ONE_LAYER_ONLY"

	// Configures the MD5 loader to not load the MD5ANIM file for
	// a MD5MESH file automatically.
	//
	// The default strategy is to look for a file with the same name but the
	// MD5ANIM extension in the same directory. If it is found, it is loaded
	// and combined with the MD5MESH file. This configuration option can be
	// used to disable this behaviour.
	//
	// Property type: bool. Default value: false.
	MD5NoAnimAutoload Prop = "IMPORT_MD5_NO_ANIM_AUTOLOAD"

	// Defines the begin of the time range for which the LWS loader
	// evaluates animations and computes aiNodeAnim's.
	//
	// Assimp provides full conversion of LightWave's envelope system, including
	// pre and post conditions. The loader computes linearly subsampled animation
	// chanels with the frame rate given in the LWS file. This property defines
	// the start time. Note: animation channels are only generated if a node
	// has at least one envelope with more tan one key assigned. This property.
	// is given in frames, '0' is the first frame. By default, if this property
	// is not set, the importer takes the animation start from the input LWS
	// file ('FirstFrame' line)<br>
	//
	// Property type: Integer. Default value: taken from file.
	//
	// @see AI_CONFIG_IMPORT_LWS_ANIM_END - end of the imported time range
	LWSAnimStart Prop = "IMPORT_LWS_ANIM_START"
	LWSAnimEnd   Prop = "IMPORT_LWS_ANIM_END"

	// Defines the output frame rate of the IRR loader.
	//
	// IRR animations are difficult to convert for Assimp and there will
	// always be a loss of quality. This setting defines how many keys per second
	// are returned by the converter.<br>
	//
	// Property type: integer. Default value: 100
	IRRAnimFPS Prop = "IMPORT_IRR_ANIM_FPS"

	// Ogre Importer will try to load this Materialfile.
	//
	// Ogre Meshes contain only the MaterialName, not the MaterialFile. If there
	// is no material file with the same name as the material, Ogre Importer will
	// try to load this file and search the material in it.
	//
	// Property type: String. Default value: guessed.
	OgreMaterialFile Prop = "IMPORT_OGRE_MATERIAL_FILE"

	// Ogre Importer detect the texture usage from its filename
	//
	// Normally, a texture is loaded as a colormap, if no target is specified in the
	// materialfile. Is this switch is enabled, texture names ending with _n, _l, _s
	// are used as normalmaps, lightmaps or specularmaps.
	//
	// Property type: Bool. Default value: false.
	OgreTexturetypeFromFilename Prop = "IMPORT_OGRE_TEXTURETYPE_FROM_FILENAME"

	// Specifies whether the IFC loader skips over IfcSpace elements.
	//
	// IfcSpace elements (and their geometric representations) are used to
	// represent, well, free space in a building storey.<br>
	//
	// Property type: Bool. Default value: true.
	IFCSkipSpaceRepresentations Prop = "IMPORT_IFC_SKIP_SPACE_REPRESENTATIONS"

	// Specifies whether the IFC loader skips over
	// shape representations of type 'Curve2D'.
	//
	// A lot of files contain both a faceted mesh representation and a outline
	// with a presentation type of 'Curve2D'. Currently Assimp doesn't convert those,
	// so turning this option off just clutters the log with errors.<br>
	//
	// Property type: Bool. Default value: true.
	IFCSkipCurveRepresentations Prop = "IMPORT_IFC_SKIP_CURVE_REPRESENTATIONS"

	// Specifies whether the IFC loader will use its own, custom triangulation
	// algorithm to triangulate wall and floor meshes.
	//
	// If this property is set to false, walls will be either triangulated by
	// #aiProcess_Triangulate or will be passed through as huge polygons with
	// faked holes (i.e. holes that are connected with the outer boundary using
	// a dummy edge). It is highly recommended to set this property to true
	// if you want triangulated data because #aiProcess_Triangulate is known to
	// have problems with the kind of polygons that the IFC loader spits out for
	// complicated meshes.
	//
	// Property type: Bool. Default value: true.
	IFCCustomTriangulation Prop = "IMPORT_IFC_CUSTOM_TRIANGULATION"

	ColladaIgnoreUpDirection Prop = "IMPORT_COLLADA_IGNORE_UP_DIRECTION"
)
