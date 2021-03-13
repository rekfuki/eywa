import moment from 'moment';

export const parseDuration = (durationStr) => {
    if (durationStr === '') {
        return null;
    }
    if (durationStr === '0') {
        // Allow 0 without a unit.
        return 0;
    }

    const durationRE = new RegExp('^(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?$');
    const matches = durationStr.match(durationRE);
    if (!matches) {
        return null;
    }

    let dur = 0;

    // Parse the match at pos `pos` in the regex and use `mult` to turn that
    // into ms, then add that value to the total parsed duration.
    const m = (pos, mult) => {
        if (matches[pos] === undefined) {
            return;
        }
        const n = parseInt(matches[pos]);
        dur += n * mult;
    };

    m(2, 1000 * 60 * 60 * 24 * 365); // y
    m(4, 1000 * 60 * 60 * 24 * 7); // w
    m(6, 1000 * 60 * 60 * 24); // d
    m(8, 1000 * 60 * 60); // h
    m(10, 1000 * 60); // m
    m(12, 1000); // s
    m(14, 1); // ms

    return dur;
};

export const formatDuration = (d) => {
    let ms = d;
    let r = '';
    if (ms === 0) {
        return '0s';
    }

    const f = (unit, mult, exact) => {
        if (exact && ms % mult !== 0) {
            return;
        }
        const v = Math.floor(ms / mult);
        if (v > 0) {
            r += `${v}${unit}`;
            ms -= v * mult;
        }
    };

    // Only format years and weeks if the remainder is zero, as it is often
    // easier to read 90d than 12w6d.
    f('y', 1000 * 60 * 60 * 24 * 365, true);
    f('w', 1000 * 60 * 60 * 24 * 7, true);

    f('d', 1000 * 60 * 60 * 24, false);
    f('h', 1000 * 60 * 60, false);
    f('m', 1000 * 60, false);
    f('s', 1000, false);
    f('ms', 1, false);

    return r;
};

export function parseTime(timeText) {
    return moment.utc(timeText).valueOf();
}

export function formatTime(time) {
    return moment.utc(time).format('yyyy-MM-DDTHH:mm:ss');
}

export function formatTimeNoSeconds(time) {
    return moment.utc(time).format('yyyy-MM-DDTHH:mm');
}

export const now = () => moment().valueOf();

export const humanizeDuration = (milliseconds) => {
    const sign = milliseconds < 0 ? '-' : '';
    const unsignedMillis = milliseconds < 0 ? -1 * milliseconds : milliseconds;
    const duration = moment.duration(unsignedMillis, 'ms');
    const ms = Math.floor(duration.milliseconds());
    const s = Math.floor(duration.seconds());
    const m = Math.floor(duration.minutes());
    const h = Math.floor(duration.hours());
    const d = Math.floor(duration.asDays());
    if (d !== 0) {
        return `${sign}${d}d ${h}h ${m}m ${s}s`;
    }
    if (h !== 0) {
        return `${sign}${h}h ${m}m ${s}s`;
    }
    if (m !== 0) {
        return `${sign}${m}m ${s}s`;
    }
    if (s !== 0) {
        return `${sign}${s}.${ms}s`;
    }
    if (unsignedMillis > 0) {
        return `${sign}${unsignedMillis.toFixed(3)}ms`;
    }
    return '0s';
};

export const getDefaultTimeRange = () => {
    const now = moment().valueOf();

    return {
        from: moment(now).subtract(6, 'hour'),
        to: now,
        raw: { from: 'now-6h', to: 'now' }
    };
}

export function roundInterval(interval) {
    switch (true) {
        // 0.015s
        case interval < 15:
            return 10; // 0.01s
        // 0.035s
        case interval < 35:
            return 20; // 0.02s
        // 0.075s
        case interval < 75:
            return 50; // 0.05s
        // 0.15s
        case interval < 150:
            return 100; // 0.1s
        // 0.35s
        case interval < 350:
            return 200; // 0.2s
        // 0.75s
        case interval < 750:
            return 500; // 0.5s
        // 1.5s
        case interval < 1500:
            return 1000; // 1s
        // 3.5s
        case interval < 3500:
            return 2000; // 2s
        // 7.5s
        case interval < 7500:
            return 5000; // 5s
        // 12.5s
        case interval < 12500:
            return 10000; // 10s
        // 17.5s
        case interval < 17500:
            return 15000; // 15s
        // 25s
        case interval < 25000:
            return 20000; // 20s
        // 45s
        case interval < 45000:
            return 30000; // 30s
        // 1.5m
        case interval < 90000:
            return 60000; // 1m
        // 3.5m
        case interval < 210000:
            return 120000; // 2m
        // 7.5m
        case interval < 450000:
            return 300000; // 5m
        // 12.5m
        case interval < 750000:
            return 600000; // 10m
        // 12.5m
        case interval < 1050000:
            return 900000; // 15m
        // 25m
        case interval < 1500000:
            return 1200000; // 20m
        // 45m
        case interval < 2700000:
            return 1800000; // 30m
        // 1.5h
        case interval < 5400000:
            return 3600000; // 1h
        // 2.5h
        case interval < 9000000:
            return 7200000; // 2h
        // 4.5h
        case interval < 16200000:
            return 10800000; // 3h
        // 9h
        case interval < 32400000:
            return 21600000; // 6h
        // 1d
        case interval < 86400000:
            return 43200000; // 12h
        // 1w
        case interval < 604800000:
            return 86400000; // 1d
        // 3w
        case interval < 1814400000:
            return 604800000; // 1w
        // 6w
        case interval < 3628800000:
            return 2592000000; // 30d
        default:
            return 31536000000; // 1y
    }
}
