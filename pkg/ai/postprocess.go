package ai

/*
#include "assimp/postprocess.h"
*/
import "C"

type PostFlags uint

const (
	// Calculates the tangents and bitangents for the imported meshes.
	//
	// Does nothing if a mesh does not have normals. You might want this post
	// processing step to be executed if you plan to use tangent space
	// calculations such as normal mapping applied to the meshes.
	//
	// The MaxSmoothingAngle configuration property allows you to specify a
	// maximum smoothing angle for the algorithm. However, usually you'll want
	// to leave it at the default value.
	CalcTangentSpace PostFlags = C.aiProcess_CalcTangentSpace

	// Identifies and joins identical vertex data sets within all imported
	// meshes.
	//
	// After this step is run, each mesh contains unique vertices, so a vertex
	// may be used by multiple faces. You usually want to use this post
	// processing step. If your application deals with indexed geometry, this
	// step is compulsory or you'll just waste rendering time.
	//
	// If this flag is not specified, no vertices are references by more than
	// one face and no index buffer is required for rendering.
	JoinIdenticalVertices PostFlags = C.aiProcess_JoinIdenticalVertices

	// Converts all the imported data to a left-handed coordinate space.
	//
	// By default the data is returned in a right-handed coordinate space
	// (which OpenGL prefers). In this space, +x points to the right, +Z points
	// towards the viewer, and +Y points upwards. In the DirectX coordinate
	// space +X points to the right, +Y points upwards, and +Z points away from
	// the viewer.
	//
	// You'll probably want to consider this flag if you use Direct3D for
	// rendering. The ConvertToLeftHanded flag supersedes this setting and
	// bundles all conversions typically required for D3D-based applications.
	MakeLeftHanded PostFlags = C.aiProcess_MakeLeftHanded

	// Triangulates all faces of all meshes.
	//
	// By default the imported mesh data might contain faces with more than 3
	// indices. For rendering you'll usually want all faces to be triangles.
	// This post processing step splits up faces with more than 3 indices into
	// triangles. Line and point primitives are not modified! If you want
	// 'triangles only' with no other kinds of primitives, try the following
	// solution:
	//  Triangleulate|SortByPType
	//  // Ignore all point and line meshes when processing data.
	Triangulate PostFlags = C.aiProcess_Triangulate

	// Removes some parts of the data structure (animations, materials, light
	// sources, cameras, textures, vertex components).
	//
	// The componets to be removed are specified as a importer property:
	// RCFlags.
	//
	// This is quite useful if you don't need all parts of the output
	// structure. Vertex colors are rarely used today for example. Calling this
	// step to remove unneeded data from the pipeline as early as possible
	// results in increased performance and a more optimized output data
	// structure.
	//
	// This step is also useful if you want to force Assimp to recompute
	// normals or tangents. The corresponding steps don't recompute them if
	// they're already there (loaded from the source asset). By using this step
	// you can make sure they are NOT there.
	//
	// This flag is a poor one, mainly because its purpose is usually
	// misunderstood. Consider the following case: a 3D model has been exported
	// from a CAD app, and it has per-face vertex colors. Vertex positions
	// can't be shared, thus JoinIdenticalVertices step fails to optimize the
	// data because of these nasty little vertex colors. Most apps don't even
	// process them, so it's all for nothing. By using this step, unneeded
	// components are excluded as early as possible thus opening more room for
	// internal optimizations.
	RemoveComponent PostFlags = C.aiProcess_RemoveComponent

	// Generates normals for all faces of all meshes.
	//
	// This is ignored if normals are already there at the time this flag is
	// evaluated. Model importers try to load them from the source file, so
	// they're usually already there. Face normals are shared between all
	// points of a single face, so a single point can have multiple normals,
	// which forces the library to duplicate vertices in some cases. The
	// JoinIdenticalVertices flag is senseless then.
	//
	// This flag cannot be specified together with GenSmoothNormals.
	GenNormals PostFlags = C.aiProcess_GenNormals

	// Generates smooth normals for all vertices in the mesh.
	//
	// This is ignored if normals are already there at the time this flag is
	// evaluated. Model importers try to load them from the source file, so
	// they're usually already there.
	//
	// The importer property GSNMaxSmoothingAngle allows you to specify an
	// angle maximum for the normal smoothing algorithm. Normals exceeding this
	// limit are not smoothed, resulting in a 'hard' seam between two faces.
	// Using a decent angle here (e.g. 80 degrees) results in very good visual
	// appearance.
	//
	// This flag cannot be specified together with GenNormals.
	GenSmoothNormals PostFlags = C.aiProcess_GenSmoothNormals

	// Splits large meshes into smaller sub-meshes.
	//
	// This is quite useful for real-time rendering, where the number of
	// triangles which can be maximally processed in a single draw-call is
	// limited by the video driver/hardware. The maximum vertex buffer is
	// usually limited too. Both requirements can be met with this step: you
	// may specify both a triangle and vertex limit for a single mesh.
	//
	// The split limits can (and should!) be set through the SLMVertexLimit and
	// SLMTriangleLimit importer properties. The default values are
	// SLMDefaultMaxVertices and SLMDefaultMaxTriangles.
	//
	// Note that splitting is generally a time-consuming task, but only if
	// there's something to split. The use of this step is recommended for most
	// users.
	SplitLargeMeshes PostFlags = C.aiProcess_SplitLargeMeshes

	// Removes the node graph and pre-transforms all vertices with
	// the local transformation matrices of their nodes.
	//
	// The output scene still contains nodes, however there is only a
	// root node with children, each one referencing only one mesh,
	// and each mesh referencing one material. For rendering, you can
	// simply render all meshes in order - you don't need to pay
	// attention to local transformations and the node hierarchy.
	// Animations are removed during this step.
	// This step is intended for applications without a scenegraph.
	// The step CAN cause some problems: if e.g. a mesh of the asset
	// contains normals and another, using the same material index, does not,
	// they will be brought together, but the first meshes's part of
	// the normal list is zeroed. However, these artifacts are rare.
	// The PTVNormalize configuration property
	// can be set to normalize the scene's spatial dimension to the -1...1
	// range.
	PreTransformVertices PostFlags = C.aiProcess_PreTransformVertices

	// Limits the number of bones simultaneously affecting a single vertex
	//  to a maximum value.
	//
	// If any vertex is affected by more than the maximum number of bones, the
	// least  important vertex weights are removed and the remaining vertex
	// weights are renormalized so that the weights still sum up to 1.
	//
	// The default bone weight limit is 4 (LBWDefaultMaxWeights), but you can
	// use the LBWMaxWeights importer property to supply your own limit to the
	// post processing step.
	//
	// If you intend to perform the skinning in hardware, this post processing
	// step might be of interest to you.
	LimitBoneWeights PostFlags = C.aiProcess_LimitBoneWeights

	// Validates the imported scene data structure.
	//
	// This makes sure that all indices are valid, all animations and bones are
	// linked correctly, all material references are correct .. etc.
	//
	// It is recommended that you capture Assimp's log output if you use this
	// flag, so you can easily find out what's wrong if a file fails the
	// validation. The validator is quite strict and will find all
	// inconsistencies in the data structure... It is recommended that plugin
	// developers use it to debug their loaders. There are two types of
	// validation failures:
	//
	// Error: There's something wrong with the imported data. Further
	//   postprocessing is not possible and the data is not usable at all.
	//   The import fails. #Importer::GetErrorString() or #aiGetErrorString()
	//   carry the error message around.
	// Warning: There are some minor issues (e.g. 1000000 animation
	//   keyframes with the same time), but further postprocessing and use
	//   of the data structure is still safe. Warning details are written
	//   to the log file, #AI_SCENE_FLAGS_VALIDATION_WARNING is set
	//   in #aiScene::mFlags
	//
	// This post-processing step is not time-consuming. Its use is not
	// compulsory, but recommended.
	ValidateDataStructure PostFlags = C.aiProcess_ValidateDataStructure

	// Reorders triangles for better vertex cache locality.
	//
	// The step tries to improve the ACMR (average post-transform vertex cache
	// miss ratio) for all meshes. The implementation runs in O(n) and is
	// roughly based on the 'tipsify' algorithm (see <a href="
	// http://www.cs.princeton.edu/gfx/pubs/Sander_2007_%3ETR/tipsy.pdf">this
	// paper</a>).
	//
	// If you intend to render huge models in hardware, this step might
	// be of interest to you. The #AI_CONFIG_PP_ICL_PTCACHE_SIZE
	// importer property can be used to fine-tune the cache optimization.
	ImproveCacheLocality PostFlags = C.aiProcess_ImproveCacheLocality

	// Searches for redundant/unreferenced materials and removes them.
	//
	// This is especially useful in combination with the
	// #aiProcess_PretransformVertices and #aiProcess_OptimizeMeshes flags.
	// Both join small meshes with equal characteristics, but they can't do
	// their work if two meshes have different materials. Because several
	// material settings are lost during Assimp's import filters,
	// (and because many exporters don't check for redundant materials), huge
	// models often have materials which are are defined several times with
	// exactly the same settings.
	//
	// Several material settings not contributing to the final appearance of
	// a surface are ignored in all comparisons (e.g. the material name).
	// So, if you're passing additional information through the
	// content pipeline (probably using // magic// material names), don't
	// specify this flag. Alternatively take a look at the
	// #AI_CONFIG_PP_RRM_EXCLUDE_LIST importer property.
	RemoveRedundantMaterials PostFlags = C.aiProcess_RemoveRedundantMaterials

	// This step tries to determine which meshes have normal vectors
	// that are facing inwards and inverts them.
	//
	// The algorithm is simple but effective:
	// the bounding box of all vertices + their normals is compared against
	// the volume of the bounding box of all vertices without their normals.
	// This works well for most objects, problems might occur with planar
	// surfaces. However, the step tries to filter such cases.
	// The step inverts all in-facing normals. Generally it is recommended
	// to enable this step, although the result is not always correct.
	FixInfacingNormals PostFlags = C.aiProcess_FixInfacingNormals

	// This step splits meshes with more than one primitive type in
	//  homogeneous sub-meshes.
	//
	//  The step is executed after the triangulation step. After the step
	//  returns, just one bit is set in aiMesh::mPrimitiveTypes. This is
	//  especially useful for real-time rendering where point and line
	//  primitives are often ignored or rendered separately.
	//  You can use the #AI_CONFIG_PP_SBP_REMOVE importer property to
	//  specify which primitive types you need. This can be used to easily
	//  exclude lines and points, which are rarely used, from the import.
	SortByPType PostFlags = C.aiProcess_SortByPType

	// This step searches all meshes for degenerate primitives and
	//  converts them to proper lines or points.
	//
	// A face is 'degenerate' if one or more of its points are identical.
	// To have the degenerate stuff not only detected and collapsed but
	// removed, try one of the following procedures:
	// <br>1. (if you support lines and points for rendering but don't
	//    want the degenerates)</br>
	//
	//   Specify the #aiProcess_FindDegenerates flag.
	//
	//   Set the #AI_CONFIG_PP_FD_REMOVE importer property to
	//       1. This will cause the step to remove degenerate triangles from the
	//       import as soon as they're detected. They won't pass any further
	//       pipeline steps.
	//
	//
	// <br>2.(if you don't support lines and points at all)</br>
	//
	//   Specify the #aiProcess_FindDegenerates flag.
	//
	//   Specify the #aiProcess_SortByPType flag. This moves line and
	//     point primitives to separate meshes.
	//
	//   Set the #AI_CONFIG_PP_SBP_REMOVE importer property to
	//       @code aiPrimitiveType_POINTS | aiPrimitiveType_LINES
	//       @endcode to cause SortByPType to reject point
	//       and line meshes from the scene.
	//
	//
	// Degenerate polygons are not necessarily evil and that's why
	// they're not removed by default. There are several file formats which
	// don't support lines or points, and some exporters bypass the
	// format specification and write them as degenerate triangles instead.
	FindDegenerates PostFlags = C.aiProcess_FindDegenerates

	// This step searches all meshes for invalid data, such as zeroed
	//  normal vectors or invalid UV coords and removes/fixes them. This is
	//  intended to get rid of some common exporter errors.
	//
	// This is especially useful for normals. If they are invalid, and
	// the step recognizes this, they will be removed and can later
	// be recomputed, i.e. by the #aiProcess_GenSmoothNormals flag.<br>
	// The step will also remove meshes that are infinitely small and reduce
	// animation tracks consisting of hundreds if redundant keys to a single
	// key. The AI_CONFIG_PP_FID_ANIM_ACCURACY config property decides
	// the accuracy of the check for duplicate animation tracks.
	FindInvalidData PostFlags = C.aiProcess_FindInvalidData

	// This step converts non-UV mappings (such as spherical or
	//  cylindrical mapping) to proper texture coordinate channels.
	//
	// Most applications will support UV mapping only, so you will
	// probably want to specify this step in every case. Note that Assimp is not
	// always able to match the original mapping implementation of the
	// 3D app which produced a model perfectly. It's always better to let the
	// modelling app compute the UV channels - 3ds max, Maya, Blender,
	// LightWave, and Modo do this for example.
	//
	// If this step is not requested, you'll need to process the
	// #AI_MATKEY_MAPPING material property in order to display all assets
	// properly.
	GenUVCoords PostFlags = C.aiProcess_GenUVCoords

	// This step applies per-texture UV transformations and bakes
	//  them into stand-alone vtexture coordinate channels.
	//
	// UV transformations are specified per-texture - see the
	// #AI_MATKEY_UVTRANSFORM material key for more information.
	// This step processes all textures with
	// transformed input UV coordinates and generates a new (pre-transformed) UV channel
	// which replaces the old channel. Most applications won't support UV
	// transformations, so you will probably want to specify this step.
	//
	// UV transformations are usually implemented in real-time apps by
	// transforming texture coordinates at vertex shader stage with a 3x3
	// (homogenous) transformation matrix.
	TransformUVCoords PostFlags = C.aiProcess_TransformUVCoords

	// This step searches for duplicate meshes and replaces them
	//  with references to the first mesh.
	//
	//  This step takes a while, so don't use it if speed is a concern.
	//  Its main purpose is to workaround the fact that many export
	//  file formats don't support instanced meshes, so exporters need to
	//  duplicate meshes. This step removes the duplicates again. Please
	//  note that Assimp does not currently support per-node material
	//  assignment to meshes, which means that identical meshes with
	//  different materials are currently // not// joined, although this is
	//  planned for future versions.
	FindInstances PostFlags = C.aiProcess_FindInstances

	// A postprocessing step to reduce the number of meshes.
	//
	//  This will, in fact, reduce the number of draw calls.
	//
	//  This is a very effective optimization and is recommended to be used
	//  together with #aiProcess_OptimizeGraph, if possible. The flag is fully
	//  compatible with both #aiProcess_SplitLargeMeshes and #aiProcess_SortByPType.
	OptimizeMeshes PostFlags = C.aiProcess_OptimizeMeshes

	// A postprocessing step to optimize the scene hierarchy.
	//
	//  Nodes without animations, bones, lights or cameras assigned are
	//  collapsed and joined.
	//
	//  Node names can be lost during this step. If you use special 'tag nodes'
	//  to pass additional information through your content pipeline, use the
	//  #AI_CONFIG_PP_OG_EXCLUDE_LIST importer property to specify a
	//  list of node names you want to be kept. Nodes matching one of the names
	//  in this list won't be touched or modified.
	//
	//  Use this flag with caution. Most simple files will be collapsed to a
	//  single node, so complex hierarchies are usually completely lost. This is not
	//  useful for editor environments, but probably a very effective
	//  optimization if you just want to get the model data, convert it to your
	//  own format, and render it as fast as possible.
	//
	//  This flag is designed to be used with #aiProcess_OptimizeMeshes for best
	//  results.
	//
	//  'Crappy' scenes with thousands of extremely small meshes packed
	//  in deeply nested nodes exist for almost all file formats.
	//  #aiProcess_OptimizeMeshes in combination with #aiProcess_OptimizeGraph
	//  usually fixes them all and makes them renderable.
	OptimizeGraph PostFlags = C.aiProcess_OptimizeGraph

	// This step flips all UV coordinates along the y-axis and adjusts
	// material settings and bitangents accordingly.
	//
	// Output UV coordinate system:
	// @code
	// 0y|0y ---------- 1x|0y
	// |                 |
	// |                 |
	// |                 |
	// 0x|1y ---------- 1x|1y
	// @endcode
	//
	// You'll probably want to consider this flag if you use Direct3D for
	// rendering. The #aiProcess_ConvertToLeftHanded flag supersedes this
	// setting and bundles all conversions typically required for D3D-based
	// applications.
	FlipUVs PostFlags = C.aiProcess_FlipUVs

	// This step adjusts the output face winding order to be CW.
	//
	// The default face winding order is counter clockwise (CCW).
	//
	// Output face order:
	// @code
	//       x2
	//
	//                         x0
	//  x1
	// @endcode
	FlipWindingOrder PostFlags = C.aiProcess_FlipWindingOrder

	// This step splits meshes with many bones into sub-meshes so that each
	// su-bmesh has fewer or as many bones as a given limit.
	SplitByBoneCount PostFlags = C.aiProcess_SplitByBoneCount

	// This step removes bones losslessly or according to some threshold.
	//
	//  In some cases (i.e. formats that require it) exporters are forced to
	//  assign dummy bone weights to otherwise static meshes assigned to
	//  animated meshes. Full, weight-based skinning is expensive while
	//  animating nodes is extremely cheap, so this step is offered to clean up
	//  the data in that regard.
	//
	//  Use #AI_CONFIG_PP_DB_THRESHOLD to control this.
	//  Use #AI_CONFIG_PP_DB_ALL_OR_NONE if you want bones removed if and
	// 	only if all bones within the scene qualify for removal.
	Debone PostFlags = C.aiProcess_Debone

	// Shortcut flag for Direct3D-based applications.
	//
	// Supersedes the MakeLeftHanded and FlipUVs and FlipWindingOrder flags.
	// The output data matches Direct3D's conventions: left-handed geometry,
	// upper-left origin for UV coordinates and finally clockwise face order,
	// suitable for CCW culling.
	ConvertToLeftHanded PostFlags = MakeLeftHanded | FlipUVs | FlipWindingOrder

	// Default postprocess configuration optimizing the data for real-time
	// rendering.
	//
	// Applications would want to use this preset to load models on end-user
	// PCs, maybe for direct use in game.
	//
	// If you're using DirectX, don't forget to combine this value with the
	// ConvertToLeftHanded step. If you don't support UV transformations in
	// your application apply the TransformUVCoords step, too.
	//
	// Note: Please take the time to read the docs for the steps enabled by
	// this preset. Some of them offer further configurable properties, while
	// some of them might not be of use for you so it might be better to not
	// specify them.
	TargetRealtimeFast PostFlags = CalcTangentSpace | GenNormals | JoinIdenticalVertices | Triangulate | GenUVCoords | SortByPType

	// Default postprocess configuration optimizing the data for real-time
	// rendering.
	//
	// Unlike TargetRealtimeFast, this configuration performs some extra
	// optimizations to improve rendering speed and to minimize memory usage.
	// It could be a good choice for a level editor environment where import
	// speed is not so important.
	//
	// If you're using DirectX, don't forget to combine this value with the
	// ConvertToLeftHanded step. If you don't support UV transformations in
	// your application apply the TransformUVCoords step, too.
	//
	// Note: Please take the time to read the docs for the steps enabled by
	// this preset. Some of them offer further configurable properties, while
	// some of them might not be of use for you so it might be better to not
	// specify them.
	TargetRealtimeQuality PostFlags = CalcTangentSpace | GenSmoothNormals | JoinIdenticalVertices | ImproveCacheLocality | LimitBoneWeights | RemoveRedundantMaterials | SplitLargeMeshes | Triangulate | GenUVCoords | SortByPType | FindDegenerates | FindInvalidData

	// Default postprocess configuration optimizing the data for real-time
	// rendering.
	//
	// This preset enables almost every optimization step to achieve perfectly
	// optimized data. It's your choice for level editor environments where
	// import speed is not important.
	//
	// If you're using DirectX, don't forget to combine this value with the
	// ConvertToLeftHanded step. If you don't support UV transformations in
	// your application apply the TransformUVCoords step, too.
	//
	// Note: Please take the time to read the docs for the steps enabled by
	// this preset. Some of them offer further configurable properties, while
	// some of them might not be of use for you so it might be better to not
	// specify them.
	TargetRealtimeMaxQuality PostFlags = FindInstances | ValidateDataStructure | OptimizeMeshes
)
