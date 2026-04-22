import { auth } from "$lib/server/auth";
import { collectionCourses } from "$lib/server/database";
import type { Actions } from "@sveltejs/kit";
import { ObjectId } from "mongodb";

export const actions: Actions = {
  getUser: async (event) => {
    const formData = await event.request.formData();
    const id = formData.get("id")?.toString() ?? "";

    const data = await auth.api.getUser({
      query: {
        id: id,
      },
      headers: event.request.headers,
    });

    return {
      user: data,
    };
  },
};
