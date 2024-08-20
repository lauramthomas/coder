# Get started quickly with Coder

## Step 1: Install Coder

### On Kubernetes (recommended for over 50 users)
<!-- I would need to know more about the product(s) and processes to determine what info should go here versus in Installation -> Kubernetes. I don't want to just copy and paste content verbatim in multiple places in the docs. Suffice it to say that information relevant to getting started quickly would go here, whereas more in-depth info would go in the the main content.  -->

### On a local machine
<!-- Same as above.  -->

### On a VM
<!-- Same as above.  -->

## Step 2: Create a template
<!-- The same as above applies here, but I took a stab at including some higher level content. -->
1. In Coder, on the **Templates** tab, click **Starter Templates**.

![Starter Templates button](../images/templates/starter-templates-button.png)

2. In **Filter**, select **Docker** and then select **Develop in Docker**.

![Choosing a starter template](../images/templates/develop-in-docker-template.png)

3. Select **Use template**.

![Using a starter template](../images/templates/use-template.png)

4. In **Create template**, enter a **Name** and a **Display name**, then scroll down
and select **Create template**.

![Creating a template](../images/templates/create-template.png)

5. When the template is ready, select **Create Workspace**.

![Create workspace](../images/templates/create-workspace.png)

6. In **New workspace**, enter a **Name** and then scroll down to select **Create Workspace**. Coder starts your new workspace from your template, and after a few seconds, your workspace is ready to use.

![New workspace](../images/templates/new-workspace.png)

![Workspace is ready](../images/templates/workspace-ready.png)

## Step 3: Create a workspace

### In the Coder UI
In your Coder instance, on the **Templates** tab, next to the template you need, select **Create Workspace**.

![Creating a workspace in the UI](./images/creating-workspace-ui.png)

When you create a workspace, you will be prompted to give it a name. You might
also be prompted to set some parameters that the template provides.

You can manage your existing templates on the **Workspaces** tab.

### From the command line

Each Coder user has their own workspaces created from [shared templates](../templates/index.md):

```shell
# create a workspace from the template; specify any variables
coder create --template="<templateName>" <workspaceName>

# show the resources behind the workspace and how to connect
coder show <workspace-name>
```

For more information about workspaces in Coder, see [Workspaces](../workspaces.md).
