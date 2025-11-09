interface Props {
  name: string;
  email: string;
  picture: string;
  onClick: () => void;
}

export const ProfileCard = ({ name, email, picture, onClick }: Props) => {
  const avatarProxyUrl = `http://localhost:8080/avatar?url=${encodeURIComponent(
    picture
  )}`;

  return (
    <div className="flex flex-col gap-4 h-auto border border-gray-200 shadow-2xs rounded-xl p-4 md:p-5 bg-transparent dark:border-neutral-700 dark:shadow-neutral-700/70">
      <div className="flex gap-4">
        <img
          src={avatarProxyUrl}
          alt={name}
          className="inline-block size-15.5 rounded-full"
        />
        <div className="">
          <h3 className="text-xl font-bold text-gray-800 dark:text-white">
            {name}
          </h3>
          <p className="mt-1 text-md font-medium text-gray-500 dark:text-neutral-200">
            {email}
          </p>
        </div>
      </div>
      <button
        className="rounded-md px-3 py-1 font-medium bg-rose-700 text-white hover:bg-rose-500"
        onClick={onClick}
      >
        Logout
      </button>
    </div>
  );
};
