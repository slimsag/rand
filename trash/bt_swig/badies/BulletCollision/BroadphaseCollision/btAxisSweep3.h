//Bullet Continuous Collision Detection and Physics Library
//Copyright (c) 2003-2006 Erwin Coumans  http://continuousphysics.com/Bullet/

//
// btAxisSweep3.h
//
// Copyright (c) 2006 Simon Hobbs
//
// This software is provided 'as-is', without any express or implied warranty. In no event will the authors be held liable for any damages arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose, including commercial applications, and to alter it and redistribute it freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not claim that you wrote the original software. If you use this software in a product, an acknowledgment in the product documentation would be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source distribution.

#ifndef BT_AXIS_SWEEP_3_H
#define BT_AXIS_SWEEP_3_H

#include "LinearMath/btVector3.h"
#include "btOverlappingPairCache.h"
#include "btBroadphaseInterface.h"
#include "btBroadphaseProxy.h"
#include "btOverlappingPairCallback.h"
#include "btDbvtBroadphase.h"

//#define DEBUG_BROADPHASE 1
#define USE_OVERLAP_TEST_ON_REMOVES 1

class btAxisSweep3InternalHandle;
class btAxisSweep3InternalEdge;

/// The internal templace class btAxisSweep3Internal implements the sweep and prune broadphase.
/// It uses quantized integers to represent the begin and end points for each of the 3 axis.
/// Dont use this class directly, use btAxisSweep3 or bt32BitAxisSweep3 instead.
template <typename BP_FP_INT_TYPE>
class btAxisSweep3Internal : public btBroadphaseInterface
{
protected:

	BP_FP_INT_TYPE	m_bpHandleMask;
	BP_FP_INT_TYPE	m_handleSentinel;

public:
	
 BT_DECLARE_ALIGNED_ALLOCATOR();
	
protected:
	btVector3 m_worldAabbMin;						// overall system bounds
	btVector3 m_worldAabbMax;						// overall system bounds

	btVector3 m_quantize;						// scaling factor for quantization

	BP_FP_INT_TYPE m_numHandles;						// number of active handles
	BP_FP_INT_TYPE m_maxHandles;						// max number of handles
	btAxisSweep3InternalHandle* m_pHandles;						// handles pool
	
	BP_FP_INT_TYPE m_firstFreeHandle;		// free handles list

	btAxisSweep3InternalEdge* m_pEdges[3];						// edge arrays for the 3 axes (each array has m_maxHandles * 2 + 2 sentinel entries)
	void* m_pEdgesRawPtr[3];

	btOverlappingPairCache* m_pairCache;

	///btOverlappingPairCallback is an additional optional user callback for adding/removing overlapping pairs, similar interface to btOverlappingPairCache.
	btOverlappingPairCallback* m_userPairCallback;
	
	bool	m_ownsPairCache;

	int	m_invalidPair;

	///additional dynamic aabb structure, used to accelerate ray cast queries.
	///can be disabled using a optional argument in the constructor
	btDbvtBroadphase*	m_raycastAccelerator;
	btOverlappingPairCache*	m_nullPairCache;


	// allocation/deallocation
	BP_FP_INT_TYPE allocHandle();
	void freeHandle(BP_FP_INT_TYPE handle);
	

	bool testOverlap2D(const btAxisSweep3InternalHandle* pHandleA, const btAxisSweep3InternalHandle* pHandleB,int axis0,int axis1);

#ifdef DEBUG_BROADPHASE
	void debugPrintAxis(int axis,bool checkCardinality=true);
#endif //DEBUG_BROADPHASE

	//Overlap* AddOverlap(BP_FP_INT_TYPE handleA, BP_FP_INT_TYPE handleB);
	//void RemoveOverlap(BP_FP_INT_TYPE handleA, BP_FP_INT_TYPE handleB);

	

	void sortMinDown(int axis, BP_FP_INT_TYPE edge, btDispatcher* dispatcher, bool updateOverlaps );
	void sortMinUp(int axis, BP_FP_INT_TYPE edge, btDispatcher* dispatcher, bool updateOverlaps );
	void sortMaxDown(int axis, BP_FP_INT_TYPE edge, btDispatcher* dispatcher, bool updateOverlaps );
	void sortMaxUp(int axis, BP_FP_INT_TYPE edge, btDispatcher* dispatcher, bool updateOverlaps );

public:

	btAxisSweep3Internal(const btVector3& worldAabbMin,const btVector3& worldAabbMax, BP_FP_INT_TYPE handleMask, BP_FP_INT_TYPE handleSentinel, BP_FP_INT_TYPE maxHandles = 16384, btOverlappingPairCache* pairCache=0,bool disableRaycastAccelerator = false);

	virtual	~btAxisSweep3Internal();

	BP_FP_INT_TYPE getNumHandles() const
	{
		return m_numHandles;
	}

	virtual void	calculateOverlappingPairs(btDispatcher* dispatcher);
	
	BP_FP_INT_TYPE addHandle(const btVector3& aabbMin,const btVector3& aabbMax, void* pOwner,short int collisionFilterGroup,short int collisionFilterMask,btDispatcher* dispatcher,void* multiSapProxy);
	void removeHandle(BP_FP_INT_TYPE handle,btDispatcher* dispatcher);
	void updateHandle(BP_FP_INT_TYPE handle, const btVector3& aabbMin,const btVector3& aabbMax,btDispatcher* dispatcher);
	SIMD_FORCE_INLINE btAxisSweep3InternalHandle* getHandle(BP_FP_INT_TYPE index) const {return m_pHandles + index;}

	virtual void resetPool(btDispatcher* dispatcher);

	void	processAllOverlappingPairs(btOverlapCallback* callback);

	//Broadphase Interface
	virtual btBroadphaseProxy*	createProxy(  const btVector3& aabbMin,  const btVector3& aabbMax,int shapeType,void* userPtr ,short int collisionFilterGroup,short int collisionFilterMask,btDispatcher* dispatcher,void* multiSapProxy);
	virtual void	destroyProxy(btBroadphaseProxy* proxy,btDispatcher* dispatcher);
	virtual void	setAabb(btBroadphaseProxy* proxy,const btVector3& aabbMin,const btVector3& aabbMax,btDispatcher* dispatcher);
	virtual void  getAabb(btBroadphaseProxy* proxy,btVector3& aabbMin, btVector3& aabbMax ) const;
	
	virtual void	rayTest(const btVector3& rayFrom,const btVector3& rayTo, btBroadphaseRayCallback& rayCallback, const btVector3& aabbMin=btVector3(0,0,0), const btVector3& aabbMax = btVector3(0,0,0));
	virtual void	aabbTest(const btVector3& aabbMin, const btVector3& aabbMax, btBroadphaseAabbCallback& callback);

	
	void quantize(BP_FP_INT_TYPE* out, const btVector3& point, int isMax) const;
	///unQuantize should be conservative: aabbMin/aabbMax should be larger then 'getAabb' result
	void unQuantize(btBroadphaseProxy* proxy,btVector3& aabbMin, btVector3& aabbMax ) const;
	
	bool	testAabbOverlap(btBroadphaseProxy* proxy0,btBroadphaseProxy* proxy1);

	btOverlappingPairCache*	getOverlappingPairCache()
	{
		return m_pairCache;
	}
	const btOverlappingPairCache*	getOverlappingPairCache() const
	{
		return m_pairCache;
	}

	void	setOverlappingPairUserCallback(btOverlappingPairCallback* pairCallback)
	{
		m_userPairCallback = pairCallback;
	}
	const btOverlappingPairCallback*	getOverlappingPairUserCallback() const
	{
		return m_userPairCallback;
	}

	///getAabb returns the axis aligned bounding box in the 'global' coordinate frame
	///will add some transform later
	virtual void getBroadphaseAabb(btVector3& aabbMin,btVector3& aabbMax) const
	{
		aabbMin = m_worldAabbMin;
		aabbMax = m_worldAabbMax;
	}

	virtual void	printStats()
	{
/*		printf("btAxisSweep3.h\n");
		printf("numHandles = %d, maxHandles = %d\n",m_numHandles,m_maxHandles);
		printf("aabbMin=%f,%f,%f,aabbMax=%f,%f,%f\n",m_worldAabbMin.getX(),m_worldAabbMin.getY(),m_worldAabbMin.getZ(),
			m_worldAabbMax.getX(),m_worldAabbMax.getY(),m_worldAabbMax.getZ());
			*/

	}

};

////////////////////////////////////////////////////////////////////






#ifdef DEBUG_BROADPHASE
#include <stdio.h>

template <typename BP_FP_INT_TYPE>
void btAxisSweep3<BP_FP_INT_TYPE>::debugPrintAxis(int axis, bool checkCardinality)
{
	int numEdges = m_pHandles[0].m_maxEdges[axis];
	printf("SAP Axis %d, numEdges=%d\n",axis,numEdges);

	int i;
	for (i=0;i<numEdges+1;i++)
	{
		btAxisSweep3InternalEdge* pEdge = m_pEdges[axis] + i;
		btAxisSweep3InternalHandle* pHandlePrev = getHandle(pEdge->m_handle);
		int handleIndex = pEdge->IsMax()? pHandlePrev->m_maxEdges[axis] : pHandlePrev->m_minEdges[axis];
		char beginOrEnd;
		beginOrEnd=pEdge->IsMax()?'E':'B';
		printf("	[%c,h=%d,p=%x,i=%d]\n",beginOrEnd,pEdge->m_handle,pEdge->m_pos,handleIndex);
	}

	if (checkCardinality)
		btAssert(numEdges == m_numHandles*2+1);
}
#endif //DEBUG_BROADPHASE


/// The btAxisSweep3 is an efficient implementation of the 3d axis sweep and prune broadphase.
/// It uses arrays rather then lists for storage of the 3 axis. Also it operates using 16 bit integer coordinates instead of floats.
/// For large worlds and many objects, use bt32BitAxisSweep3 or btDbvtBroadphase instead. bt32BitAxisSweep3 has higher precision and allows more then 16384 objects at the cost of more memory and bit of performance.
class btAxisSweep3 : public btAxisSweep3Internal<unsigned short int>
{
public:

	btAxisSweep3(const btVector3& worldAabbMin,const btVector3& worldAabbMax, unsigned short int maxHandles = 16384, btOverlappingPairCache* pairCache = 0, bool disableRaycastAccelerator = false);

};

/// The bt32BitAxisSweep3 allows higher precision quantization and more objects compared to the btAxisSweep3 sweep and prune.
/// This comes at the cost of more memory per handle, and a bit slower performance.
/// It uses arrays rather then lists for storage of the 3 axis.
class bt32BitAxisSweep3 : public btAxisSweep3Internal<unsigned int>
{
public:

	bt32BitAxisSweep3(const btVector3& worldAabbMin,const btVector3& worldAabbMax, unsigned int maxHandles = 1500000, btOverlappingPairCache* pairCache = 0, bool disableRaycastAccelerator = false);

};

#endif



	class btAxisSweep3InternalEdge
	{
	public:
		BP_FP_INT_TYPE m_pos;			// low bit is min/max
		BP_FP_INT_TYPE m_handle;

		BP_FP_INT_TYPE IsMax() const {return static_cast<BP_FP_INT_TYPE>(m_pos & 1);}
	};

	class	btAxisSweep3InternalHandle : public btBroadphaseProxy
	{
	public:
	BT_DECLARE_ALIGNED_ALLOCATOR();
	
		// indexes into the edge arrays
		BP_FP_INT_TYPE m_minEdges[3], m_maxEdges[3];		// 6 * 2 = 12
//		BP_FP_INT_TYPE m_uniqueId;
		btBroadphaseProxy*	m_dbvtProxy;//for faster raycast
		//void* m_pOwner; this is now in btBroadphaseProxy.m_clientObject
	
		SIMD_FORCE_INLINE void SetNextFree(BP_FP_INT_TYPE next) {m_minEdges[0] = next;}
		SIMD_FORCE_INLINE BP_FP_INT_TYPE GetNextFree() const {return m_minEdges[0];}
	};		// 24 bytes + 24 for Edge structures = 44 bytes total per entry

