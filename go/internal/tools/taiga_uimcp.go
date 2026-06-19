output, cloneErr := cmd.CombinedOutput()
if cloneErr != nil {
    return err(fmt.Sprintf("git clone failed: %s\n%s", cloneErr.Error(), string(output)))
}

entries, readErr := os.ReadDir(componentsDir)
if readErr != nil {
    // Fallback: try listing the repo root
    entries, readErr = os.ReadDir(repoPath)
    if readErr != nil {
        return err(fmt.Sprintf("failed to read directory: %s", readErr.Error()))

}
}