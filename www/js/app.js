function get_actions(info) {
    actions = [];
    if (info.installed_version) {
        actions.push('open');
        if (info.current_version != info.installed_version) {
            actions.push('upgrade');
        }
        actions.push('remove');
    } else {
        actions.push('install');
    }

    return actions;
}