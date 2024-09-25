
# beehive
ReMenu Monorepo

------------------------

# Table of content

## [Quick start with project](#quick-start)
* [Clone project in your local system](#clone-project-in-your-local-system)
* [Syncing with main branch in repository](#syncing-with-main-branch-in-repository)
* 



------------------------

# Quick Start

  ## Clone project in your local system

* You can first specify username and email locally or globally:
```bash
# to define them in current project
git config user.name <username>
git config user.email <email>

# to define them for all projects globally
git config --global user.name <username>
git config --global user.email <email>
```

* Before cloning repository, you can use this command to save your username & password:
```bash
git config --global credential.helper store
```

* Clone repo with https method:
```bash
git clone https://git.gocasts.ir/remenu/beehive.git
```

* To see all branches: `git branch -a`
* To create new branch in your local system: 
```bash
git checkout -b <username/feature/task-name>

#e.g.
git checkout -b mehdi/basket/add-item
```
* After adding your code, you must add it stage: `git add .`
* Then you must commit your changes: `git commit -m "Your Explanation About Your Changes"`
* Push your changes to remote repository
  **Note**: If it's first time to push new branch, you can use `-u` flage
```bash
# for first time
git push -u origin <username/feature/task-name>

# in normal cases
git push origin mehdi/basket/add-item
```
* Now you can see a message to suggest you to send a **Pull Request** in your remote repository
* Put your comments and send your pull request


  ## Syncing with main branch in repository

If main repository updated, you need to update your local branches. You can do these steps in the following:

* If at first, you want to see changes with your branches, you can `fetch` first. 
* Then with `git status` or `git diff`, you can see the changes
```bash
git fetch

git status

# or
git diff origin/main
```

* Now you need to merge it with your desire branch (I prefer to checkout to main branch and merge it)
```bash
git checkout main

git merge origin/main
```
  * **Note**: Sometimes when you want to merge, you will face with conflict. So you need to solve it and then add new commit before merging

* Sometimes you want to fetch and merge together, so you must use `pull` command:
```bash
git pull origin main
```
* Now you can checkout to your desire branch and coding again  
