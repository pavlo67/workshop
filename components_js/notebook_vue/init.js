import b    from '../basis';
import {ep} from '../swagger_convertor';

let cfg = {};

function init(data) {
    if (!(data instanceof Object)) {
        return;
    }

    if ('backend' in data) {
        // TODO: do it safely!!!

        cfg.readEp   = ep(data.backend, "read").replace("/{id}", "");
        cfg.saveEp   = ep(data.backend, "save");
        cfg.removeEp = ep(data.backend, "remove").replace("/{id}", "");

        cfg.recentEp = ep(data.backend, "recent");
        cfg.tagsEp   = ep(data.backend, "tags");
        cfg.taggedEp = ep(data.backend, "tagged");
    }

    if ('eventBus' in data) {
        cfg.eventBus = data.eventBus;
        cfg.eventBus.$on('user', user => {

            console.log(333333333);
            cfg.eventBus.$on('user', user => {
                console.log(444444444);
                cfg.user = user;
            });
        });
    }

    if ('vue' in data) {
        cfg.vue = data.vue;
    }

}

export { cfg, init };
