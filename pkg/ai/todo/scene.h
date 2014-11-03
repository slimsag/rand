/*
---------------------------------------------------------------------------
Open Asset Import Library (assimp)
---------------------------------------------------------------------------

Copyright (c) 2006-2012, assimp team

All rights reserved.

Redistribution and use of this software in source and binary forms, 
with or without modification, are permitted provided that the following 
conditions are met:

* Redistributions of source code must retain the above
  copyright notice, this list of conditions and the
  following disclaimer.

* Redistributions in binary form must reproduce the above
  copyright notice, this list of conditions and the
  following disclaimer in the documentation and/or other
  materials provided with the distribution.

* Neither the name of the assimp team, nor the names of its
  contributors may be used to endorse or promote products
  derived from this software without specific prior
  written permission of the assimp team.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS 
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT 
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT 
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT 
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY 
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT 
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE 
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
---------------------------------------------------------------------------
*/

/** @file aiScene.h
 *  @brief Defines the data structures in which the imported scene is returned.
 */
#ifndef __AI_SCENE_H_INC__
#define __AI_SCENE_H_INC__

#include "types.h"
#include "texture.h"
#include "mesh.h"
#include "light.h"
#include "camera.h"
#include "material.h"
#include "anim.h"
#include "metadata.h"

#ifdef __cplusplus
extern "C" {
#endif


struct aiNode
{
	/** The name of the node. 
	 *
	 * The name might be empty (length of zero) but all nodes which 
	 * need to be referenced by either bones or animations are named.
	 * Multiple nodes may have the same name, except for nodes which are referenced
	 * by bones (see #aiBone and #aiMesh::mBones). Their names *must* be unique.
	 * 
	 * Cameras and lights reference a specific node by name - if there
	 * are multiple nodes with this name, they are assigned to each of them.
	 * <br>
	 * There are no limitations with regard to the characters contained in
	 * the name string as it is usually taken directly from the source file. 
	 * 
	 * Implementations should be able to handle tokens such as whitespace, tabs,
	 * line feeds, quotation marks, ampersands etc.
	 *
	 * Sometimes assimp introduces new nodes not present in the source file
	 * into the hierarchy (usually out of necessity because sometimes the
	 * source hierarchy format is simply not compatible). Their names are
	 * surrounded by @verbatim <> @endverbatim e.g.
	 *  @verbatim<DummyRootNode> @endverbatim.
	 */
	C_STRUCT aiString mName;

	/** The transformation relative to the node's parent. */
	C_STRUCT aiMatrix4x4 mTransformation;

	/** Parent node. NULL if this node is the root node. */
	C_STRUCT aiNode* mParent;

	/** The number of child nodes of this node. */
	unsigned int mNumChildren;

	/** The child nodes of this node. NULL if mNumChildren is 0. */
	C_STRUCT aiNode** mChildren;

	/** The number of meshes of this node. */
	unsigned int mNumMeshes;

	/** The meshes of this node. Each entry is an index into the mesh */
	unsigned int* mMeshes;

	/** Metadata associated with this node or NULL if there is no metadata.
	  *  Whether any metadata is generated depends on the source file format. See the
	  * @link importer_notes @endlink page for more information on every source file
	  * format. Importers that don't document any metadata don't write any. 
	  */
	C_STRUCT aiMetadata* mMetaData;

#ifdef __cplusplus
	/** Searches for a node with a specific name, beginning at this
	 *  nodes. Normally you will call this method on the root node
	 *  of the scene.
	 * 
	 *  @param name Name to search for
	 *  @return NULL or a valid Node if the search was successful.
	 */
	inline const aiNode* FindNode(const aiString& name) const
	{
		return FindNode(name.data);
	}


	inline aiNode* FindNode(const aiString& name)
	{
		return FindNode(name.data);
	}


	/** @override
	 */
	inline const aiNode* FindNode(const char* name) const
	{
		if (!::strcmp( mName.data,name))return this;
		for (unsigned int i = 0; i < mNumChildren;++i)
		{
			const aiNode* const p = mChildren[i]->FindNode(name);
			if (p) {
				return p;
			}
		}
		// there is definitely no sub-node with this name
		return NULL;
	}

	inline aiNode* FindNode(const char* name) 
	{
		if (!::strcmp( mName.data,name))return this;
		for (unsigned int i = 0; i < mNumChildren;++i)
		{
			aiNode* const p = mChildren[i]->FindNode(name);
			if (p) {
				return p;
			}
		}
		// there is definitely no sub-node with this name
		return NULL;
	}

#endif // __cplusplus
};



