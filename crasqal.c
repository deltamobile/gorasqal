#include <rasqal.h>
#include "crasqal.h"
#include "_cgo_export.h"

void gorasqal_set_log_handler(rasqal_world *world, void *user_data) {
    rasqal_world_set_log_handler(world, user_data, (raptor_log_handler)GoRasqal_log_handler);
}

raptor_sequence *goraptor_new_sequence() {
    return raptor_new_sequence((raptor_data_free_handler)rasqal_free_data_graph, NULL);
}
