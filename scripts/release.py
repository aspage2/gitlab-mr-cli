
import json
import requests
import click
import os
from urllib.parse import urljoin
from pathlib import Path

# Normally https://gitlab.com/
GITLAB_HOST = os.environ["CI_PROJECT_URL"]

ACCESS_TOKEN = os.environ["PRIVATE_TOKEN"]
PROJECT_ID = os.environ["CI_PROJECT_ID"]
RELEASE_TAG = os.environ["CI_COMMIT_TAG"]


def upload_files(filenames, dry_run):
    
    upload_url = urljoin(
        GITLAB_HOST,
        f"/api/v4/projects/{PROJECT_ID}/uploads",
    )

    ret = {}
    
    for fname in filenames:
        with open(fname, 'rb') as f:
            if not dry_run:
                resp = requests.post(
                    upload_url,
                    headers={"Authorization": f"Bearer {ACCESS_TOKEN}"},
                    files={"file": f},
                )
                resp.raise_for_status()
                data = resp.json()
                print(data)
                url = data["full_path"]
                alt = data["alt"]
            else:
                alt = fname
                url = f"/data/{alt}"
            ret[alt] = url
    return ret


    
@click.command()
@click.option("--dry-run/--no-dry-run", default=False) 
@click.argument("upload_file", nargs=-1)
def main(upload_file, dry_run):

    release_url = urljoin(
        GITLAB_HOST,
        f"/api/v4/projects/{PROJECT_ID}/releases",
    )

    req_body = {"tag_name": RELEASE_TAG}
    
    if len(upload_file) != 0:
        uploaded = upload_files(upload_file, dry_run)
        req_body["assets"] = {
            "links": [
                {
                    "name": Path(fname).name,
                    "url": urljoin(GITLAB_HOST, fullpath)
                }
                for fname, fullpath in uploaded.items()
            ]
        }

    if not dry_run:
        resp = requests.post(
            release_url,
            headers={"Authorization": f"Bearer {ACCESS_TOKEN}"},
            json=req_body,
        )
        resp.raise_for_status()
    else:
        print(json.dumps(req_body, indent=2))
    print("Done")


if __name__ == "__main__":
    main()
