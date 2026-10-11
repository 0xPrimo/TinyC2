#ifndef _DEBUG_H_
#define _DEBUG_H_

#if defined( _DEBUG )
#include <stdio.h>
#define DBG_PRINTF( format, ... )                                                                                                                    \
    {                                                                                                                                                \
        printf( "[DEBUG::%s::%d] " format, __FUNCTION__, __LINE__, ##__VA_ARGS__ );                                                                  \
    }
#else
#define DBG_PRINTF( format, ... )                                                                                                                    \
    {                                                                                                                                                \
        ;                                                                                                                                            \
    }
#endif

#endif // !_DEBUG_H_
