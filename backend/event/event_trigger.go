package event

import (
	"github.com/syncloud/platform/snap"
	"log"
	"os/exec"
)

type EventTrigger struct {
	snap *snap.Snap
}

const (
	SNAP_INFO_CMD = "snap info"
	SNAP_RUN_CMD  = "snap run"
)

func New() *EventTrigger {
	return &EventTrigger{
		snap: snap.New(),
	}
}

func (storage *EventTrigger) Trigger(event string) error {
	log.Println("Running: ", SNAP_INFO_CMD, event)
	out, err := exec.Command(SNAP_INFO_CMD, "testapp").CombinedOutput()
	if err != nil {
		log.Printf("snap info failed: %v", err)
		return err
	}
	log.Printf("snap info output: %v", out)
	return nil
	/*
		def _trigger_app_event(self, action):
		        for app in self.installer.installed_all_apps():
		            app_id = app.app.id
		            try:
		                info = check_output('snap info {0}'.format(app_id), shell=True).decode()
		                command_name = '{0}.{1}'.format(app_id, action)
		                if command_name in info:
		                    command = 'snap run {0}'.format(command_name)
		                    self.log.info('executing {0}'.format(command))
		                    output = check_output(command, shell=True, stderr=STDOUT).decode()
		                    print(output)
		            except CalledProcessError as e:
		                self.log.error('event output {0}'.format(e.output.decode()))
		                if e.stderr:
		                    self.log.error('event error {0}'.format(e.stderr.decode()))
		                if e.stdout:
		                    self.log.error('event stdout {0}'.format(e.stdout.decode()))
		                self.log.error(traceback.format_exc())
	*/
}
