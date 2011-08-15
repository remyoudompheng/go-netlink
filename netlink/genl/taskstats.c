#include <linux/taskstats.h>

typedef struct taskstats $TaskStats;

enum {
     $TASKSTATS_VERSION = TASKSTATS_VERSION,
     $TS_COMM_LEN = TS_COMM_LEN,
};

enum {
	$TASKSTATS_CMD_UNSPEC = 0,	/* Reserved */
	$TASKSTATS_CMD_GET,		/* user->kernel request/get-response */
	$TASKSTATS_CMD_NEW,		/* kernel->user event */
	$__TASKSTATS_CMD_MAX,
};

enum {
	$TASKSTATS_TYPE_UNSPEC = 0,	/* Reserved */
	$TASKSTATS_TYPE_PID,		/* Process id */
	$TASKSTATS_TYPE_TGID,		/* Thread group id */
	$TASKSTATS_TYPE_STATS,		/* taskstats structure */
	$TASKSTATS_TYPE_AGGR_PID,	/* contains pid + stats */
	$TASKSTATS_TYPE_AGGR_TGID,	/* contains tgid + stats */
	$TASKSTATS_TYPE_NULL,		/* contains nothing */
	$__TASKSTATS_TYPE_MAX,
};

enum {
	$TASKSTATS_CMD_ATTR_UNSPEC = 0,
	$TASKSTATS_CMD_ATTR_PID,
	$TASKSTATS_CMD_ATTR_TGID,
	$TASKSTATS_CMD_ATTR_REGISTER_CPUMASK,
	$TASKSTATS_CMD_ATTR_DEREGISTER_CPUMASK,
	$__TASKSTATS_CMD_ATTR_MAX,
};

