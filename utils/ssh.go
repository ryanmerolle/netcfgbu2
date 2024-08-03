package utils

import (
	"regexp"
	"time"

	"github.com/ryanmerolle/netcfgbu2/models"

	"golang.org/x/crypto/ssh"
)

// appendDefaultPort appends ":22" to the host if it does not end with a colon followed by a valid integer
func appendDefaultPort(host string) string {
	re := regexp.MustCompile(`:([0-9]+)$`)
	if !re.MatchString(host) {
		Log.Debugf("Appending default port to host: %s", host)
		return host + ":22"
	}
	return host
}

func RunSSHCommand(device models.Device, username, password, command string, timeout int) (string, error) {
	//Log.Debugf("Running SSH command on host: %s, user: %s, command: %s, timeout: %d", device.Host, username, command, timeout)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(timeout) * time.Second,
	}

	address := appendDefaultPort(device.Host)

	conn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		Log.Warnf("ERROR: Failed to dial SSH: %v", err)
		return "", err
	}

	defer func() {
		Log.Debugf("Closing SSH connection to %s", address)
		conn.Close()
	}()

	Log.Infof("LOGIN: %s (%s) timeout=%ds as %s", device.Hostname, device.Platform, timeout, username)
	session, err := conn.NewSession()
	if err != nil {
		Log.Warnf("ERROR: Failed to create SSH session: %v", err)
		return "", err
	}
	Log.Infof("CONNECTED: %s", device.Hostname)

	defer func() {
		Log.Infof("CLOSED: %s", device.Hostname)
		session.Close()
	}()

	//Log.Debugf("Executing command: %s", command)
	output, err := session.CombinedOutput(command)
	if err != nil {
		Log.Warnf("ERROR: Failed to execute command: %v", err)
		return "", err
	}
	Log.Infof("GET-CONFIG: %s (%s) timeout=%ds", device.Hostname, device.Platform, timeout)

	// Ensure the output has a single newline at the end
	outputStr := string(output)
	outputStr = ensureSingleNewline(outputStr)

	Log.Debugf("Command output: %s", outputStr)
	return outputStr, nil
}
